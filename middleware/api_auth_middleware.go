package middleware

import (
	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/services/cache"
)

func ApiAuth(ctx context.Context) {
	token := getUserToken(ctx)
	userToken := cache.UserTokenCache.Get(token)

	// 没找到授权
	if userToken == nil || userToken.Status == model.UserTokenStatusDisabled {
		notLogin(ctx)
		return
	}
	// 授权过期
	if userToken.ExpiredAt <= simple.NowTimestamp() {
		notLogin(ctx)
		return
	}

	ctx.Next()
}

// 从请求体中获取UserToken
func getUserToken(ctx context.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

func notLogin(ctx context.Context) {
	_, _ = ctx.JSON(simple.Error(simple.ErrorNotLogin))
	ctx.StopExecution()
}
