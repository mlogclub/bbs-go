package api

import (
	"bbs-go/controllers/render"
	"bbs-go/pkg/qq"
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type QQLoginController struct {
	Ctx iris.Context
}

// 获取QQ登录授权地址
func (c *QQLoginController) GetAuthorize() *web.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.QQ {
		return web.JsonErrorMsg("QQ登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	url := qq.AuthorizeUrl(map[string]string{"ref": ref})
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取QQ回调信息获取
func (c *QQLoginController) GetCallback() *web.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.QQ {
		return web.JsonErrorMsg("QQ登录/注册已禁用")
	}

	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByQQ(code, state)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return web.JsonError(codeErr)
	} else {
		return render.BuildLoginSuccess(user, "")
	}
}
