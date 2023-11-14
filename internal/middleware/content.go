package middleware

import (
	"github.com/valyala/fasthttp"
)

func ContentTypeHtmlMiddleware(ctx *fasthttp.RequestCtx) error {
	ctx.SetContentType("text/html; charset=utf-8")
	return nil
}

