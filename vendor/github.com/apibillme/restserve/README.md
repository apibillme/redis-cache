# restserve - Fast REST server for Go (golang)

[![Go Report](https://goreportcard.com/badge/github.com/apibillme/restserve)](https://goreportcard.com/report/github.com/apibillme/restserve) [![GolangCI](https://golangci.com/badges/github.com/apibillme/restserve.svg)](https://golangci.com/r/github.com/apibillme/restserve) [![Travis](https://travis-ci.org/apibillme/restserve.svg?branch=master)](https://travis-ci.org/apibillme/restserve#) [![codecov](https://codecov.io/gh/apibillme/restserve/branch/master/graph/badge.svg)](https://codecov.io/gh/apibillme/restserve) ![License](https://img.shields.io/github/license/apibillme/restserve.svg) ![Maintenance](https://img.shields.io/maintenance/yes/2018.svg) [![GoDoc](https://godoc.org/github.com/apibillme/restserve?status.svg)](https://godoc.org/github.com/apibillme/restserve)


This fast REST server uses [fasthttp](https://github.com/valyala/fasthttp#readme) for performance.

## Features:
* Nice middleware
* RESTful
* About 10x faster than using net/http
* About 200 LOC

```go
import (
    "github.com/apibillme/restserve"
    "github.com/valyala/fasthttp"
)
```

```go
app := restserve.New()

app.Finally func(err error, ctx *fasthttp.RequestCtx)

app.Use(router string, func(ctx *fasthttp.RequestCtx, next func(error)) // register non-error-middleware

app.Use(router string, func(err error, ctx *fasthttp.RequestCtx, next func(error))) // register error-middleware

app.Listen(port string) // listen a port. app.Listen(":8080")

app.Get(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http GET method

app.Post(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http POST method

app.Put(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http PUT method

app.Head(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http HEAD method

app.Patch(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http PATCH method

app.Delete(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http DELETE method

app.Options(router string, func(ctx *fasthttp.RequestCtx, next func(error))) // only handle http OPTIONS method
```

## Example App

* [restserve-demo](https://github.com/apibillme/restserve-demo)
