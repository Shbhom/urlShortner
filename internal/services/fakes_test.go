package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	appredis "github.com/shbhom/urlShortner/internal/db/redis"
	"github.com/shbhom/urlShortner/internal/models"
)

type FakeURLRepo struct {
	sync.RWMutex
	urls        map[string]models.UrlData
	seq         uint64
	bulkUpdates []map[string]string
}

func NewFakeURLRepo() *FakeURLRepo {
	return &FakeURLRepo{
		urls:        make(map[string]models.UrlData),
		bulkUpdates: make([]map[string]string, 0),
	}
}

func (f *FakeURLRepo) GetNextSequence(ctx context.Context) (uint64, error) {
	f.Lock()
	defer f.Unlock()
	f.seq++
	return f.seq, nil
}

func (f *FakeURLRepo) AddUrl(ctx context.Context, data models.UrlData) error {
	f.Lock()
	defer f.Unlock()
	f.urls[data.ShortCode] = data
	return nil
}

func (f *FakeURLRepo) GetUrlByCode(ctx context.Context, short_code string) (string, error) {
	f.RLock()
	defer f.RUnlock()
	if data, ok := f.urls[short_code]; ok {
		return data.TargetUrl, nil
	}
	return "", errors.New("url not found")
}

func (f *FakeURLRepo) GetBulkUpdates() []map[string]string {
	f.RLock()
	defer f.RUnlock()
	return f.bulkUpdates
}

func (f *FakeURLRepo) BulkUpdateUrlLastInvokation(ctx context.Context, data map[string]string) error {
	f.Lock()
	defer f.Unlock()
	f.bulkUpdates = append(f.bulkUpdates, data)
	return nil
}

type FakeCacheRepo struct {
	sync.RWMutex
	store  map[string]string
	hashes map[string]map[string]string
}

func NewFakeCacheRepo() *FakeCacheRepo {
	return &FakeCacheRepo{
		store:  make(map[string]string),
		hashes: make(map[string]map[string]string),
	}
}

func (f *FakeCacheRepo) Get(ctx context.Context, shortCode string) (string, error) {
	f.RLock()
	defer f.RUnlock()
	if val, ok := f.store[shortCode]; ok {
		return val, nil
	}
	return "", redis.Nil
}

func (f *FakeCacheRepo) Set(ctx context.Context, data models.UrlData) error {
	f.Lock()
	defer f.Unlock()
	f.store[data.ShortCode] = data.TargetUrl
	return nil
}

func (f *FakeCacheRepo) Rename(ctx context.Context, oldKey, newKey string) error {
	f.Lock()
	defer f.Unlock()

	found := false
	if val, ok := f.store[oldKey]; ok {
		f.store[newKey] = val
		delete(f.store, oldKey)
		found = true
	}
	if hash, ok := f.hashes[oldKey]; ok {
		f.hashes[newKey] = hash
		delete(f.hashes, oldKey)
		found = true
	}

	if !found {
		return errors.New("ERR no such key")
	}
	return nil
}

func (f *FakeCacheRepo) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	f.RLock()
	defer f.RUnlock()
	if hash, ok := f.hashes[key]; ok {
		ret := make(map[string]string)
		for k, v := range hash {
			ret[k] = v
		}
		return ret, nil
	}
	return make(map[string]string), nil
}

func (f *FakeCacheRepo) Delete(ctx context.Context, key string) error {
	f.Lock()
	defer f.Unlock()
	delete(f.store, key)
	delete(f.hashes, key)
	return nil
}

func (f *FakeCacheRepo) RecordInvokation(ctx context.Context, code string) error {
	f.Lock()
	defer f.Unlock()
	key := appredis.ANALYTICS_KEY
	if f.hashes[key] == nil {
		f.hashes[key] = make(map[string]string)
	}
	f.hashes[key][code] = fmt.Sprintf("%d", time.Now().Unix())
	return nil
}
