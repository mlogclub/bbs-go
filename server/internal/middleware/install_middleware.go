package middleware

import (
	"bbs-go/internal/pkg/config"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

func Install(ctx iris.Context) {
	installed := config.Instance.Installed
	if installed {
		ctx.Next()
		return
	}

	path := ctx.Path()
	if strings.HasPrefix(path, "/api/install/") || path == "/api/config/configs" || path == "/api/user/current" {
		ctx.Next()
		return
	}

	_ = ctx.JSON(web.JsonErrorCode(-1, "Please install first"))
	ctx.StopExecution()
}
