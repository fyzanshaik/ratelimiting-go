# Token Bucket Rate Limiter

## How It Works

Token bucket algorithm: each client gets a bucket with tokens. Every request consumes 1 token. Tokens refill at a constant rate. No tokens = rate limited.

## Layers

1. **TokenBucket** - Core algorithm implementation
2. **RateLimiter** - Manages buckets per client (by IP)
3. **Middleware** - HTTP wrapper that checks limits before handling requests

## Features & Behavior

- **Client Identification**: correctly identifies clients by IP address, stripping ports and respecting `X-Forwarded-For` / `X-Real-IP` headers for proxy support.
- **Memory Management**: includes an automatic cleanup process. Buckets that have been inactive for more than 5 minutes are removed every minute to prevent memory leaks.

## Usage

```bash
go mod init ratelimiter-example
go run main.go
```

Test it:

```bash
for i in {1..15}; do curl http://localhost:8080/; done
```

After 10 requests you'll get rate limited (429 error). Tokens refill at 2/second.