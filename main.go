package main

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/cast"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/gammazero/workerpool"
	"github.com/valyala/fasthttp"

	"github.com/apibillme/cache"
	"github.com/apibillme/restserve"
)

// for stubbing
var redisGet = redisGetImpl
var redisNewClient = redis.NewClient
var cachedGet = cachedGetImpl
var cachedSet = cachedSetImpl

func cachedSetImpl(cached cache.Cache, key string, value string) bool {
	return cached.Set(key, value)
}

func cachedGetImpl(cached cache.Cache, key string) (interface{}, bool) {
	return cached.Get(key)
}

func redisGetImpl(client *redis.Client, key string) string {
	return client.Get(key).Val()
}

func getRedisKey(address string, cached cache.Cache, key string) (string, error) {
	// setup client
	client := redisNewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// check cache first - if found return value
	cachedVal, ok := cachedGet(cached, key)
	// if not in cache then get value and save in cache - return err if not
	if !ok {
		backingVal := redisGet(client, key)
		ok = cachedSet(cached, key, backingVal)
		if !ok {
			return "", errors.New("cannot set new value in cache")
		}
		return backingVal, nil
	}
	return cast.ToString(cachedVal), nil
}

func server(keyCapacity int, globalTTL time.Duration, address string, port string) {
	app := restserve.New()
	wp := workerpool.New(100)
	cached := cache.New(keyCapacity, cache.WithTTL(globalTTL*time.Second))

	app.Get("/", func(ctx *fasthttp.RequestCtx, next func(error)) {
		key := cast.ToString(ctx.QueryArgs().Peek("key"))
		wp.SubmitWait(func() {
			val, err := getRedisKey(address, cached, key)
			if err != nil {
				ctx.SetStatusCode(500)
				ctx.SetBodyString(`{"error":"` + err.Error() + `"}`)
				ctx.SetContentType("application/json")
			} else {
				ctx.SetStatusCode(200)
				ctx.SetBodyString(`{"value": "` + val + `"}`)
				ctx.SetContentType("application/json")
			}
		})
	})

	app.Listen(":" + port)
}

func main() {
	pflag.String("address", "localhost:6379", "redis backing server address")
	pflag.Int("ttl", 360, "global ttl for redis cache in seconds")
	pflag.Int("keycap", 128, "key capacity for redis cache")
	pflag.String("port", "8000", "port for redis cache")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	globalTTL := time.Duration(viper.GetInt("ttl"))
	keyCapacity := cast.ToInt(viper.Get("keycap"))
	address := cast.ToString(viper.Get("address"))
	port := cast.ToString(viper.Get("port"))

	server(keyCapacity, globalTTL, address, port)
}
