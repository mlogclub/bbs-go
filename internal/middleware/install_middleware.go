package middleware

import (
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/ginx"
	"net/http"
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

	// Keep the installation-required signal separate from ordinary business
	// errors. The SPA uses this status to enter the install flow; treating every
	// response with errorCode -1 as an install requirement can cause a redirect
	// loop after the site has been installed.
	ginx.WriteHttpStatusJSON(ctx, http.StatusPreconditionRequired, ginx.ErrorCode(-1, "Please install first"))
	ctx.Abort()
}
