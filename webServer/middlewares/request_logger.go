package middlewares

import (
	"github.com/techrail/ground/logger"
	"github.com/valyala/fasthttp"
)

// RequestLogger is supposed to log the info about the hit received much the same way nginx logs its requests
func RequestLogger(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// TODO: Fill this one with the request
		logger.Println("I#1N30SD - ")
		logger.LogWithContext(ctx, "L#1MVZRU - Hit the CheckOpLogRequest Middleware")
		handler(ctx)
	}
}
