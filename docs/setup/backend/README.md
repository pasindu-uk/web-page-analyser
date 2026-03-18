# Backend Setup Guide

Step-by-step instructions for setting up and running the Go backend locally.

## Prerequisites

- **Go 1.26 or later** — [Download](https://go.dev/dl/)
- **MySQL 8+** (optional) — only needed if you want to persist analysis history

Verify Go is installed:

```bash
go version
# go version go1.26.1 darwin/arm64
```

## Step 1: Clone the repository

```bash
git clone https://github.com/pasindu-uk/web-page-analyser.git
cd web-page-analyser/backend
```

## Step 2: Install dependencies

Go modules are used for dependency management. Dependencies are downloaded automatically on build, but you can fetch them explicitly:

```bash
go mod download
```

This pulls two external dependencies:
- `golang.org/x/net` — HTML parser for DOM traversal
- `github.com/go-sql-driver/mysql` — MySQL driver (only used when `MYSQL_DSN` is set)

## Step 3: Configure environment variables

Copy the example config file:

```bash
cp .env.example .env.local
```

Edit `.env.local` with your preferred settings:

```env
PORT=8080
REQUEST_TIMEOUT=10s
MAX_LINK_CHECK_WORKERS=5
LOG_LEVEL=info
MYSQL_DSN=
```

The backend automatically reads `.env.local` first, then falls back to `.env`. Environment variables set in your shell always take precedence over file values.

### Configuration options

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port the HTTP server listens on |
| `REQUEST_TIMEOUT` | `10s` | Timeout for fetching remote pages (Go duration format, e.g., `10s`, `30s`) |
| `MAX_LINK_CHECK_WORKERS` | `5` | Number of goroutines in the link checker worker pool |
| `LOG_LEVEL` | `info` | Logging level: `debug`, `info`, `warn`, `error` |
| `MYSQL_DSN` | _(empty)_ | MySQL connection string (see Step 5) |

## Step 4: Run the backend

```bash
go run ./cmd/api
```

You should see:

```
2026/03/16 10:00:00 INFO server starting addr=:8080
```

### Verify it's working

```bash
# Health check
curl http://localhost:8080/health
# ok

# Analyze a URL
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

## Step 5: Set up MySQL (optional)

MySQL is **optional**. Without it, the app works perfectly — you just won't have analysis history stored between restarts.

### 5a. Start MySQL

If you don't have MySQL running locally, you can use Docker:

```bash
docker run -d \
  --name web-analyzer-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=web_analyzer \
  -p 3306:3306 \
  mysql:8
```

### 5b. Set the connection string

Update your `.env.local`:

```env
MYSQL_DSN=root:root@tcp(127.0.0.1:3306)/web_analyzer?parseTime=true
```

The DSN format is: `user:password@tcp(host:port)/database?parseTime=true`

### 5c. Migrations

Migrations run **automatically** when the backend starts with `MYSQL_DSN` configured. The SQL files are embedded in the binary — no manual migration step is needed.

You should see this in the logs:

```
2026/03/16 10:00:00 INFO MySQL persistence enabled
```

### 5d. Verify history endpoint

The history endpoint uses an **in-memory cache** — the first call queries MySQL, and subsequent calls return cached results until a new analysis is saved. After analyzing a few URLs, check that history is stored:

```bash
curl http://localhost:8080/api/analyses
```

### 5e. Clear the cache

You can manually flush the in-memory cache without restarting the server:

```bash
curl -X DELETE http://localhost:8080/api/cache
# {"status":"cache cleared"}
```

The next call to `GET /api/analyses` will re-query MySQL. The frontend also provides a **Clear Cache** button in the History section.

## Step 6: Run tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/analyzer/...

# Run a single test
go test -run TestAnalyze_HTML5Doctype ./internal/analyzer/...

# Run with verbose output
go test -v ./...
```

## Step 7: Build a binary

```bash
go build -o bin/api ./cmd/api
./bin/api
```

The binary is self-contained — migration SQL files are embedded in it via Go's `//go:embed` directive.

## Project structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go               # Entrypoint — wires everything together
├── internal/
│   ├── analyzer/
│   │   ├── analyzer.go           # HTML parsing (doctype, title, headings, login form)
│   │   └── analyzer_test.go
│   ├── config/
│   │   └── config.go             # Loads env vars and .env files
│   ├── fetcher/
│   │   ├── fetcher.go            # HTTP client for fetching remote pages
│   │   └── fetcher_test.go
│   ├── handler/
│   │   ├── handler.go            # HTTP handlers (POST /api/analyze, GET /api/analyses, DELETE /api/cache)
│   │   └── handler_test.go
│   ├── logger/
│   │   └── logger.go             # Configures structured logging (slog)
│   ├── model/
│   │   └── model.go              # Shared request/response types
│   ├── repository/
│   │   ├── repository.go         # MySQL repository (Save/List)
│   │   ├── cache.go              # In-memory cache (decorator over Repository)
│   │   ├── cache_test.go
│   │   ├── migrations.go         # Embedded SQL migration runner
│   │   ├── migrations/
│   │   │   └── 001_create_analyses.sql
│   │   └── repository_test.go
│   └── service/
│       ├── service.go            # Business logic: fetch → analyze → check links → save
│       └── service_test.go
├── .env.example
└── go.mod
```

## Troubleshooting

### Port already in use

```
Error: listen tcp :8080: bind: address already in use
```

Change the port in `.env.local`:
```env
PORT=9090
```

### MySQL connection refused

Make sure MySQL is running and the DSN is correct:
```bash
mysql -u root -p -h 127.0.0.1 -P 3306
```

### Slow link checking

Increase the worker pool size:
```env
MAX_LINK_CHECK_WORKERS=10
```
