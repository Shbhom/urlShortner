package metrics

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

type mockDriver struct{}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	return &mockConn{}, nil
}

type mockConn struct{}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (c *mockConn) Close() error {
	return nil
}

func (c *mockConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func init() {
	sql.Register("mock_driver", &mockDriver{})
}

func TestDBStatsCollector(t *testing.T) {
	db, err := sql.Open("mock_driver", "")
	if err != nil {
		t.Fatalf("Failed to open mock db: %v", err)
	}
	defer db.Close()

	collector := NewDBStatsCollector(db)

	// Test Describe
	chDesc := make(chan *prometheus.Desc, 10)
	go func() {
		collector.Describe(chDesc)
		close(chDesc)
	}()

	descCount := 0
	for desc := range chDesc {
		if desc == nil {
			t.Error("Received nil descriptor")
		}
		descCount++
	}
	if descCount != 5 {
		t.Errorf("Expected 5 descriptors, got %d", descCount)
	}

	// Test Collect
	chMetric := make(chan prometheus.Metric, 10)
	go func() {
		collector.Collect(chMetric)
		close(chMetric)
	}()

	metricCount := 0
	for metric := range chMetric {
		if metric == nil {
			t.Error("Received nil metric")
		}
		metricCount++
	}
	if metricCount != 5 {
		t.Errorf("Expected 5 metrics, got %d", metricCount)
	}
}
