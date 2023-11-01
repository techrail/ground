package middlewares

import (
	"github.com/techrail/ground/constants/customCtxKey"
	"github.com/techrail/ground/logger"
	types "github.com/techrail/ground/typs"
	"github.com/valyala/fasthttp"
)

func SetRandomVar(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger.Println("D#1MT2BV - Hit the SetRandomVar Middleware!")

		// Check that the header is present or not.
		ctx.SetUserValue(customCtxKey.RandomValue, types.GetRandomAlphaString(20))
		handler(ctx)
	}
}
