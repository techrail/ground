package render

import "github.com/valyala/fasthttp"

func addFixedHeaders(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderServer, "Apache 2.4")
	ctx.Response.Header.Set(fasthttp.HeaderXPoweredBy, "PHP/7.2.12")
}
