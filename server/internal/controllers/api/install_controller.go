package api

import (
	"bbs-go/internal/install"
	"bbs-go/internal/pkg/config"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type InstallController struct {
	Ctx iris.Context
}

// 安装状态
func (c *InstallController) GetStatus() *web.JsonResult {
	cfg := config.Instance
	return web.JsonData(map[string]any{
		"installed": cfg.Installed,
	})
}

func (c *InstallController) PostTest_db_connection() *web.JsonResult {
	var req install.DbConfigReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return web.JsonError(err)
	}

	// 检查是否已安装
	if config.Instance.Installed {
		return web.JsonErrorMsg("系统已安装，无需重复安装")
	}

	// // 尝试连接数据库
	// sqlDB, err := sql.Open("mysql", req.GetConnStr())
	// if err != nil {
	// 	return web.JsonError(err)
	// }
	// defer sqlDB.Close()

	// if err = sqlDB.Ping(); err != nil {
	// 	return web.JsonError(err)
	// }

	if err := install.TestDbConnection(req); err != nil {
		return web.JsonError(err)
	}

	return web.JsonSuccess()
}

func (c *InstallController) PostInstall() *web.JsonResult {
	var req install.InstallReq
	if err := c.Ctx.ReadJSON(&req); err != nil {
		return web.JsonError(err)
	}

	// 检查是否已安装
	if config.Instance.Installed {
		return web.JsonErrorMsg("系统已安装，无需重复安装")
	}

	// 初始化数据库
	if err := install.Install(req); err != nil {
		return web.JsonErrorMsg("数据库初始化失败：" + err.Error())
	}

	return web.JsonData(true)
}
