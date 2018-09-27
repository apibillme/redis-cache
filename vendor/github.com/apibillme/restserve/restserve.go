package restserve

import (
	"github.com/valyala/fasthttp"
)

type (
	handleFunc func(*fasthttp.RequestCtx, func(error))
	errorsFunc func(error, *fasthttp.RequestCtx, func(error))
	midware    struct {
		route      []byte
		index      int
		middleWare handleFunc
		errorsWare errorsFunc
	}
	stack         []midware
	finalResponse func(err error, ctx *fasthttp.RequestCtx)
)

func (stack *stack) push(m midware) {
	*stack = append(*stack, m)
}

// Gorestserve is a light-weight web middleware framework
type Gorestserve struct {
	// non-error-middleware stack
	// store middleware when call app.Use()/Get/Post etc.
	stack stack

	// error-middleware stack
	// store error-middleware when call app.Use()/Get/Post etc.
	errStack stack

	// when the last error-middleware or non-error-middleware
	// calls  next(nil) or next(error), there is no more middleware.
	// In this case (see func finalRes()):
	// 1) if Finally == nil, rest will do below:
	//	  * ctx.ResetBody()
	//	  * ctx.SetStatusCode(fasthttp.StatusNotFound)
	//	  * ctx.SetBodyString("Not Found")
	// 2) if Finally is rewrote, how to response to clients by your code.
	Finally finalResponse
}

// New a Gorestserve instance
func New() *Gorestserve {
	return &Gorestserve{}
}

// Use is use to Register a middleware
// The first parameter of app.Use we call it router,
// and the second parameter we call it middleware,
// all the middlewares will be pushed to stacks.
//
// There are two stack arrays to store the midllewares which are accepted
// by app.Use. One is used to store non-error-middleware
// and the other one is used to store error-middleware
//
// func(ctx *fasthttp.RequestCtx, next func(error)),
// we call it non-error-middleware
//
// http request will execute each middleware one-by-one until a middleware
// does not call next(nil) within it.
//
// func(err error, ctx *fasthttp.RequestCtx, next func(error)), we call it
// error-middleware, only execute by call next(error) within a middleware
func (rest *Gorestserve) Use(route string, handler interface{}) {
	if route == "" || route[0] != '/' {
		panic("The first params of Use func must be a string which start with '/'")
	}

	if route == "/" {
		route = ""
	}

	midHandle, mOk := handler.(func(*fasthttp.RequestCtx, func(error)))
	errHandle, eOk := handler.(func(error, *fasthttp.RequestCtx, func(error)))

	if !mOk && !eOk {
		panic("The second params of Use func must be a" +
			"\n\tfunc(*fasthttp.RequestCtx, func(error)) or" +
			"\n\tfunc(error, *fasthttp.RequestCtx, func(error)) type")
	}

	if mOk {
		rest.stack.push(midware{
			route:      []byte(route),
			index:      len(rest.errStack),
			middleWare: midHandle})
		return
	}

	if eOk {
		rest.errStack.push(midware{
			route:      []byte(route),
			index:      len(rest.stack),
			errorsWare: errHandle})
		return
	}
}

// Post register a middleware only handle POST method
func (rest *Gorestserve) Post(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("POST"), handler)
}

// Put register a middleware only handle PUT method
func (rest *Gorestserve) Put(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("PUT"), handler)
}

// Get register a middleware only handle GET method
func (rest *Gorestserve) Get(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("GET"), handler)
}

// Delete register a middleware only handle DELETE method
func (rest *Gorestserve) Delete(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("DELETE"), handler)
}

// Options register a middleware only handle OPTIONS method
func (rest *Gorestserve) Options(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("OPTIONS"), handler)
}

// Patch register a middleware only handle PATCH method
func (rest *Gorestserve) Patch(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("PATCH"), handler)
}

// Head register a middleware only handle HEAD method
func (rest *Gorestserve) Head(route string, handler handleFunc) {
	rest.httpMethod(route, []byte("HEAD"), handler)
}

// Listen is used to listen a port
func (rest *Gorestserve) Listen(port string) error {
	return fasthttp.ListenAndServe(port, rest.Handler)
}

// Handler an incoming http request , rest will compare ctx.Path() with []byte(router),
// The compare rules are below: (see handle())
//	* if ctx.Path() equals to []byte(router) , it matches.
//	* if ctx.Path() starts with []byte(router), and len(ctx.Path()) > len(router),
//	  and ctx.Path()[len(router)] is '/' or '?', it matches.
//	* if the router is "/"，it means this router matches any http request
func (rest *Gorestserve) Handler(ctx *fasthttp.RequestCtx) {
	var (
		index    = 0
		errIndex = 0
		nxt      func(error)
		err      error
	)
	nxt = func(err error) {
		var m midware

		if err != nil {
			if errIndex >= len(rest.errStack) {
				rest.finalRes(err, ctx)
				return
			}
			m = rest.errStack[errIndex]
			errIndex++
			index = m.index
		} else {
			if index >= len(rest.stack) {
				rest.finalRes(nil, ctx)
				return
			}
			m = rest.stack[index]
			index++
			errIndex = m.index
		}

		handle(err, ctx, m, nxt)
	}

	nxt(err)
}

// The finalResponse to clients
// see Finally
func (rest *Gorestserve) finalRes(err error, ctx *fasthttp.RequestCtx) {
	if rest.Finally != nil {
		rest.Finally(err, ctx)
		return
	}
	ctx.ResetBody()
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetBodyString("Not Found")
}

// Call rest.Use(),
// only handle the given method
func (rest *Gorestserve) httpMethod(route string, method []byte, handler handleFunc) {
	rest.Use(route, func(ctx *fasthttp.RequestCtx, next func(error)) {
		if sliceCompare(ctx.Method(), method) {
			handler(ctx, next)
			return
		}
		next(nil)
	})
}

func sliceCompare(src, dest []byte) bool {
	if len(src) != len(dest) {
		return false
	}

	return sliceDiff(src, dest)
}

func sliceContains(src, dest []byte) bool {
	if len(src) < len(dest) {
		return false
	}

	return sliceDiff(src, dest)
}

func sliceDiff(src, dest []byte) bool {
	for i, w := range dest {
		if src[i] != w {
			return false
		}
	}
	return true
}

// see Handler()
func handle(err error, ctx *fasthttp.RequestCtx, m midware, n func(error)) {
	url := ctx.Path()
	urlLen := len(url)
	rouLen := len(m.route)

	if !sliceContains(url, m.route) {
		n(err)
		return
	}

	if urlLen > rouLen && url[rouLen] != '/' && url[rouLen] != '?' {
		n(err)
		return
	}

	if err != nil {
		m.errorsWare(err, ctx, n)
		return
	}

	m.middleWare(ctx, n)
	return
}
