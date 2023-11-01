package middlewares

import (
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/logger"
	"github.com/techrail/ground/uuid"
	"github.com/valyala/fasthttp"
)

func SetRequestId(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Println("D#1MR7SH- Hit the SetRequestId Middleware")

		// TODO: Implement the logging to the request context
		// Check that the header is present or not.
		requestId := ctx.Request.Header.Peek(customHeaders.RequestId)

		if requestId != nil {
			ctx.SetUserValue(customCtxKey.RequestId, string(requestId))
		} else {
			ctx.SetUserValue(customCtxKey.RequestId, uuid.GetNewUlid().String())
		}

		handler(ctx)
	}
}
