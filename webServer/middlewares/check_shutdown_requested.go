package middlewares

import (
	"github.com/techrail/ground/core"
	"github.com/techrail/ground/logger"
	"github.com/valyala/fasthttp"
)

// CheckShutdownRequested checks if the web server has been requested to be shut down.
// If it is supposed to be shut down, then a response is sent and request is not served
// Otherwise, the request is served as usual
func CheckShutdownRequested(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Println("I#1MTZVK - Hit the CheckShutdownRequested Middleware")
		//logger.LogInfoWithContext(ctx, "L#1MTZVK - Hit the CheckShutdownRequested Middleware")
		//if appRuntime.ShutdownRequested() {
		if core.State().WebServerShutdownRequested.Load() {
			ctx.Response.SetBodyString(`{"message":"Web server is shutting down and is not accepting new requests."`)
			// NOTE: once we have the rendering functions, we should replace the above statement with something like below

			//render.JsonWithFailure(
			//	ctx, fasthttp.StatusTeapot,
			//	"W#ZOMBIE",
			//	"Server is in process of shutting down",
			//	"")
			return
		}
		handler(ctx)
	}
}
