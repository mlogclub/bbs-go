package api

import (
	"bbs-go/internal/install"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/locales"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"
)

// 安装状态
func InstallStatus(ctx *gin.Context) {

	cfg := config.Instance
	ginx.WriteJSON(ctx, map[string]any{
		"installed":          cfg.Installed,
		"dockerBuiltinMysql": install.IsDockerBuiltinMySQLInstall(),
		"dbType":             cfg.DB.Type,
	})

}

func InstallTestDbConnection(ctx *gin.Context) {
	var req install.DbConfigReq
	if err := ginx.BindJSON(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	// 检查是否已安装
	if config.Instance.Installed {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("install.already_installed")))
		return
	}

	if err := install.TestDbConnection(req); err != nil {
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
