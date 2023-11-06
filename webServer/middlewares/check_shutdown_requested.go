package middlewares

import (
	"github.com/techrail/ground/core"
	"github.com/techrail/ground/logger"
	"github.com/techrail/ground/render"
	"github.com/valyala/fasthttp"
)

// CheckShutdownRequested checks if the web server has been requested to be shut down.
// If it is supposed to be shut down, then a response is sent and request is not served
// Otherwise, the request is served as usual
func CheckShutdownRequested(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.LogWithContext(ctx, "D#1MTZVK - Hit the CheckShutdownRequested Middleware")
		if core.State().WebServerShutdownRequested.Load() {
			ctx.Response.SetBodyString(`{"message":"Web server is shutting down and is not accepting new requests."`)

			render.JsonWithFailure(
				ctx, fasthttp.StatusTeapot,
				"W#ZOMBIE",
				"Server is in process of shutting down",
				"")
			return
		}
		handler(ctx)
	}
}
