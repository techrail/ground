package middlewares

import (
	"fmt"
	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/constants/customHeaders"
	"github.com/techrail/ground/logger"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// CheckOpLogRequest sets the optional variable for including Operational Log in the context depending
// on the supplied header. This value can be helpful in debugging
func CheckOpLogRequest(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		// Todo: contextual logging here
		logger.LogWithContext(ctx, "L#1MVZRU - Hit the CheckOpLogRequest Middleware")
		// Check that the header is present or not.
		requestingOpLog := ctx.Request.Header.Peek(customHeaders.OpLogRequestValue)

		if requestingOpLog != nil {
			// Check if the value matches the expected header value
			if strings.ToUpper(string(requestingOpLog)) == strings.ToUpper(constants.OpLogRequestValue) {
				ctx.SetUserValue(customCtxKey.OpLogRequested, true)
				ctx.SetUserValue(customCtxKey.CtxOperationLogContent, []string{
					fmt.Sprintf("L#1MVZRY - CheckOpLogRequest Middleware execution at %v", time.Now().UTC()),
				})
			}
		}
		handler(ctx)
	}
}
