package webServer

import "github.com/valyala/fasthttp"

type Middleware func(handler fasthttp.RequestHandler) fasthttp.RequestHandler
