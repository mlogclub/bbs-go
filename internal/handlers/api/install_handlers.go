package api

import (
	"bbs-go/internal/install"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/locales"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
)

// 安装状态
func InstallStatus(ctx *gin.Context) {

	cfg := config.Instance
	ginx.WriteJSON(ctx, map[string]any{
		"installed":           cfg.Installed,
		"dockerBuiltinMysql":  install.IsDockerBuiltinMySQLInstall(),
		"dockerBuiltinDbType": install.DockerBuiltinDBType(),
		"dbType":              cfg.DB.Type,
	})

}

func InstallTestDbConnection(ctx *gin.Context) {
	var req install.DbConfigReq
	if err := ginx.BindJSON(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if config.Instance.Installed {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("install.already_installed")))
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	if err := install.TestDbConnection(timeoutCtx, req); err != nil {
		// 判断是否因超时取消
		if errors.Is(err, context.DeadlineExceeded) {
			ginx.WriteJSON(ctx, ginx.ErrorMessage("database connection timeout, please check your network or firewall"))
			return
		}
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, nil)
}

func InstallInstall(ctx *gin.Context) {
	var req install.InstallReq
	if err := ginx.BindJSON(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	// 检查是否已安装
	if config.Instance.Installed {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("install.already_installed")))
		return
	}

	// 初始化数据库
	if err := install.Install(req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, true)

}
