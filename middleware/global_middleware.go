package middleware

import (
	"github.com/kataras/iris/context"

	"github.com/mlogclub/mlog/controllers/render"
)

type GlobalMiddleware struct {
}

func NewGlobalMiddleware() context.Handler {
	middleware := &GlobalMiddleware{}
	return middleware.GlobalMiddlewareHandle
}

func (m *GlobalMiddleware) GlobalMiddlewareHandle(ctx context.Context) {
	ctx.ViewData("CurrentUser", render.BuildCurrentUser(ctx))
	ctx.Next()
}
