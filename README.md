# High-Performance URL Shortener

A highly scalable, production-ready URL shortener built in Go. It uses PostgreSQL for persistent storage and Redis as a read-through cache and high-throughput buffer.

## 🚀 Overview

This service is designed to be extremely fast and resilient under high load. It strictly separates the read and write paths, ensuring that viral links do not overwhelm the database with write-locks.

## 🧠 Key Design Decisions

### 1. Short Code Generation
We use the `sqids-go` library coupled with a PostgreSQL `SEQUENCE` (`nextval('url_counter_seq')`).
*   **Why:** This mathematically guarantees 100% collision-free, short, and URL-safe codes natively. It completely eliminates the need for expensive database "collision checks" or while-loop retries during URL creation.

### 2. Access Tracking & Data Purging
To keep the database lean over time, we need to track the `last_invokation` timestamp of each URL so we can safely purge old/abandoned links.
*   **The Problem:** Updating Postgres on every single URL read would cause massive write-contention and lockups for viral URLs.
*   **The Solution:** Read hits are pushed into a non-blocking Go channel and buffered in a Redis Hash. A background worker periodically (every 60 seconds) flushes these buffered timestamps to Postgres in bulk.
*   **Why:** Even if a specific shortened link is clicked 10,000 times in a single minute, it results in exactly **1 database write** to update the timestamp. This keeps the read-path ultra-fast while fully enabling safe data purging.

### 3. Read-Through Caching
Redis acts as the primary read layer.
*   **Why:** Cache misses gracefully fall back to Postgres and automatically populate the cache. Subsequent reads bypass the database entirely, protecting Postgres from read spikes.

### 4. Graceful Shutdown
The server manages worker lifecycles using multiple `sync.WaitGroup`s.
*   **Why:** When the container receives a `SIGTERM` (e.g., during a deployment), it guarantees that all in-flight read-hits in the Go channel and Redis buffer are safely flushed to Postgres before the process exits, ensuring zero data loss.

## ⚙️ Environment Variables

Create a `config_prod.json` (or `.env` equivalents) with the following variables:

### Required
| Variable | Description |
| :--- | :--- |
| `DB_URL` | PostgreSQL connection string |
| `REDIS_ADDR` | Redis host address and port (e.g., `localhost:6379`) |
| `API_PORT` | Port for the HTTP server to listen on |
| `BASE_URL` | The domain used for generating the shortened links |

### Optional
| Variable | Default | Description |
| :--- | :--- | :--- |
| `URL_TTL` | `60` | Cache time-to-live in minutes |
| `SHORT_CODE_MIN_LEN` | `6` | Minimum length of the generated alias |
| `CHAIN_PATH` | `""` | Path to the TLS certificate chain (only if API_PORT=443) |
| `PEM_PATH` | `""` | Path to the TLS private key (only if API_PORT=443) |

## 🛠️ How to Run

### Local Development

1. Ensure Postgres and Redis are running locally.
2. Setup your configuration file in `~/.shortner/config_local.json` or export the environment variables.
3. Run the application:
```bash
go run cmd/server/main.go
```

### Docker (Production)

This repository includes a multi-stage `Dockerfile` which builds a highly optimized and secure Alpine-based Go image.

```bash
# Build the image
docker build -t url-shortener .

# Run the container
docker run -p 8080:8080 \
  -e DB_URL="postgres://user:pass@host:5432/db" \
  -e REDIS_ADDR="redis:6379" \
  -e API_PORT=8080 \
  -e BASE_URL="https://sho.rt" \
  url-shortener
```