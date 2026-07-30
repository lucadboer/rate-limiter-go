# rate-limiter-go

HTTP rate limiting middleware for Go services, keyed by **API token when present and by client IP
otherwise**, with counter state kept outside the process so it holds across replicas.

## The problem

An in-memory rate limiter breaks the moment you run more than one instance — each replica keeps its
own count, so the effective limit becomes `limit × replicas`. This keeps counters in Redis behind a
narrow interface, so the limit is global and the backing store is replaceable.

## How the limit is chosen

The middleware reads the `API_KEY` request header:

- **Header present** → limit by token, using `RATE_LIMIT_TOKEN`.
- **Header absent** → limit by IP, using `RATE_LIMIT_IP`.

Token limits deliberately win over IP limits, so an authenticated caller behind a shared NAT is not
throttled by its neighbours. Over the limit returns **HTTP 429**.

## Stack

Go · go-redis/v8 · Redis · testify · Docker Compose

## Running it

```bash
docker compose up -d     # Redis on 6379
go run .                 # server on :8080
```

Configuration comes from `.env` (committed — local values only):

| Variable | Meaning |
|---|---|
| `REDIS_HOST`, `REDIS_PORT` | Redis address |
| `RATE_LIMIT_IP` | requests allowed per IP per window |
| `RATE_LIMIT_TOKEN` | requests allowed per token per window |
| `BLOCK_DURATION` | window length in seconds |

Trip the IP limit with the committed defaults (10 per 300s):

```bash
for i in $(seq 1 12); do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/; done
# ten 200s, then 429s

curl -H 'API_KEY: my-token' http://localhost:8080/     # counted against the token budget instead
```

## Tests

```bash
docker compose up -d     # the tests talk to a real Redis on localhost:6379
go test ./...
```

`limiter/limiter_test.go` covers the IP and token paths. It calls `FlushAll` on setup and sleeps past
the window to assert the counter resets, so a full run takes ~20s.

## Technical decisions worth noting

**`Store` is a two-method interface** (`Incr`, `Expire`) declared in `limiter/store.go`. `RedisStore`
is the only implementation; swapping in Memcached or an in-memory store for tests means implementing
two methods and changing one line in `main.go`. The limiter itself never imports a Redis package.

**The window is set by the first request, not renewed per hit.** `Incr` returns the new count, and
only when it comes back `1` does the limiter call `Expire`. That makes it a fixed window: the TTL is
attached once and the counter dies with it, so a caller cannot extend their own block by continuing
to hammer the endpoint.

**Store errors fail closed.** If Redis is unreachable, `isRateLimited` returns `true` and the request
is rejected. The tradeoff is deliberate — a Redis outage degrades into refusing traffic rather than
silently removing every limit.

## Notes

Course challenge from the Full Cycle Go post-graduate program, kept public as a code sample.

## License

MIT
