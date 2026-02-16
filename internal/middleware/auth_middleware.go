package middleware

import (
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/services"

	"github.com/kataras/iris/v12"
)

func AuthMiddleware(ctx iris.Context) {
	if config.Instance.Installed {
		if user := services.UserTokenService.GetCurrent(ctx); user != nil {
			common.SetCurrentUser(ctx, user)
		}
	}
	ctx.Next()
}
