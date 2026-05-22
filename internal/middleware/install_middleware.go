package middleware

import (
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/ginx"
	"strings"

	"github.com/gin-gonic/gin"
)

func InstallMiddleware(ctx *gin.Context) {
	if config.Instance.Installed {
		ctx.Next()
		return
	}

	path := ctx.Request.URL.Path
	if strings.HasPrefix(path, "/api/install/") || path == "/api/config/configs" || path == "/api/user/current" {
		ctx.Next()
		return
	}

	ginx.WriteJSON(ctx, ginx.ErrorCode(-1, "Please install first"))
	ctx.Abort()
}
