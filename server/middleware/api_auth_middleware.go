package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services/cache"
)

// AdminAuth 后台权限
func AdminAuth(ctx iris.Context) {
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

	user := cache.UserCache.Get(userToken.UserId)
	userInfo := render.BuildUser(user)
	if userInfo == nil || !userInfo.HasRole("管理员") {
		_, _ = ctx.JSON(simple.JsonErrorCode(2, "无权限"))
		ctx.StopExecution()
		return
	}

	ctx.Next()
}

// 从请求体中获取UserToken
func getUserToken(ctx iris.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

func notLogin(ctx iris.Context) {
	_, _ = ctx.JSON(simple.JsonError(simple.ErrorNotLogin))
	ctx.StopExecution()
}
