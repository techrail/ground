package middlewares

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/uuid"
	"github.com/valyala/fasthttp"
	"math/rand"
	"time"
)

func SetRequestId(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Check for the request ID in the request
		requestId := ctx.Request.Header.Peek("X-Request-ID")
		if len(requestId) == 0 {
			// and if it is not supplied, then create one using ULID
			entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
			ms := ulid.Timestamp(time.Now())
			fmt.Println(ulid.New(ms, entropy))
			// 01G65Z755AFWAKHE12NY0CQ9FH
		}

		// Print
		// TODO: Implement the logging to the request context
		fmt.Println("I#1MR7SH- Hit the SetRequestId Middleware")
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
