# goratelimiter

A fixed-window rate limiter built in Go, backed by [goredis](https://github.com/nomen06/goredis) — a Redis-compatible server I built from scratch.

## How it works

Each client is identified by IP address. On every request:

1. `INCR client:<ip>` — atomically increments the request count for that client
2. If this is the first request in the window, `EXPIRE client:<ip> <window>` starts a fresh time window
3. If the count exceeds the configured limit, the request is rejected with `429 Too Many Requests`
4. Otherwise, the request is allowed through

```
Request
   ↓
Extract client IP
   ↓
INCR client:<ip>
   ↓
First request? → EXPIRE client:<ip> <window>
   ↓
count > limit?  →  yes → 429 Too Many Requests
   ↓ no
Allow request
```

## Run it

You need goredis running first:

```bash
git clone https://github.com/nomen06/goredis
cd goredis
go run *.go
```

Then in a separate terminal, run the rate limiter:

```bash
git clone https://github.com/nomen06/goratelimiter
cd goratelimiter
go run ./limiter
```

Server starts on `localhost:8080`.

## Test it

```bash
curl http://localhost:8080/
curl http://localhost:8080/
curl http://localhost:8080/
curl http://localhost:8080/   # rejected once the limit is hit
```

## Benchmark

Tested with [hey](https://github.com/rakyll/hey) — 1000 requests, 50 concurrent clients:

```
Requests/sec: 12,449
Average latency: 3.9ms
```

Handles over **12,000 requests/sec** with sub-5ms average latency under concurrent load.

## Project structure

```
limiter/
  limiter.go      — Limiter struct, Allow() logic, HTTP handler, server entrypoint
redis/
  client.go     — TCP client speaking RESP to talk to goredis
```

## Why this design

Built on a fixed-window algorithm using INCR and EXPIRE because both are atomic operations in goredis, avoiding race conditions when multiple requests arrive simultaneously for the same client. Storing rate-limit state in goredis instead of in-memory means the limiter is one step closer to working across multiple server instances rather than just inside one process.