# redis-cache - HTTP JSON caching proxy for Redis

[![Go Report](https://goreportcard.com/badge/github.com/apibillme/redis-cache)](https://goreportcard.com/report/github.com/apibillme/redis-cache) [![GolangCI](https://golangci.com/badges/github.com/apibillme/redis-cache.svg)](https://golangci.com/r/github.com/apibillme/redis-cache) [![Travis](https://travis-ci.org/apibillme/redis-cache.svg?branch=master)](https://travis-ci.org/apibillme/redis-cache#) [![codecov](https://codecov.io/gh/apibillme/redis-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/apibillme/redis-cache) ![License](https://img.shields.io/github/license/apibillme/redis-cache.svg) ![Maintenance](https://img.shields.io/maintenance/yes/2018.svg) [![GoDoc](https://godoc.org/github.com/apibillme/redis-cache?status.svg)](https://godoc.org/github.com/apibillme/redis-cache)

## Architecture Overview (what the code does)

- The proxy webserver listens on `/redis` for GET requests with the param `key`.
- Up to 100 client requests at once will run using parallel concurrent processing while the rest will be queued. The client connection remains open until their request is processed.
- The client request key/value for the GET request is cached using a fixed capacity, global expiry, LRU eviction cache.
- Less than 100 LOC.

## Make and Run
- note: flags are optional

```bash
make test
./redis-cache --ttl 360 --keycap 128 --port 8000 --address localhost:6379
```

## Docker
```bash
docker-compose up
```

## Complexity of Caching - Big-O
[O(1)](https://github.com/apibillme/cache#benchmark)


## Example Client

```go
import(
    "github.com/apibillme/restly"
)

func main() {
    req := restly.New()
    res, statusCode, err := restly.GetJSON(req, "http://localhost:8000/redis", `?key=123`)
    if err != nil {
        if statusCode == 200 {
            value := res.Get("value").String()
            log.Println(value)
        } else {
            value := res.Get("error").String()
            log.Println(value)
        }
    }
}
```

## Load Testing
[redis-cache-benchmark](https://github.com/apibillme/redis-cache-benchmark)

The first key set expires after one second and the test ran for 5 seconds - run against docker-compose in this repo on local Mac.

```
Total Requests: 5000
Requests Per Second: 1000
Success Ratio: 100%
Max: 183.256182ms
Mean: 15.357077ms
50th percentile: 730.348µs
95th percentile: 99.38549ms
99th percentile: 149.668733ms
```

## Not Implemented
- Redis RPC - took too long to write the [parser](https://github.com/apibillme/redis-rpc) and dealing with `net` package and parallel concurrent processing was something I didn't solve.
