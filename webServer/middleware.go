package webServer

import "github.com/valyala/fasthttp"

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler

type MiddlewareSet []Middleware

func chain(h fasthttp.RequestHandler, m ...Middleware) fasthttp.RequestHandler {
	// if our chain is done, use the original handler
	if len(m) == 0 {
		return h
	}
	// otherwise nest the handler functions
	return m[0](chain(h, m[1:len(m)]...))
}
