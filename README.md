#redis-cache - provides an HTTP JSON caching proxy for Redis

## Architecture Overview (what the code does)

> The proxy webserver listens on the base URL ("/") for GET requests with the param ("key").
> Up to 100 client requests at once will run using parallel concurrent processing while the rest will be queued. The client connection remains open until their request is processed.
> The client request key/value for the GET request is cached using a fixed capacity, global expiry, LRU eviction cache.
> Less than 100 LOC.

## Make and Run
- note: flags are optional

```bash
make test
./redis-cache --ttl 360 --keycap 128 --port 8000 --address localhost:6379
```

## Example

```go
import(
    "github.com/apibillme/restly"
)

func main() {
    req := restly.New()
    res, statusCode, err := restly.GetJSON(req, "http://localhost:8000/", `?key=123`)
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

## Not Implemented
- Redis RPC - took too long to write the parser and dealing with `net` and parallel concurrent processing was something I didn't solve.
