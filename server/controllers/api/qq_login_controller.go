package api

import (
	"server/controllers/render"
	"server/pkg/qq"
	"server/services"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type QQLoginController struct {
	Ctx iris.Context
}

// 获取QQ登录授权地址
func (c *QQLoginController) GetAuthorize() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	param := make(map[string]string)
	if user != nil {
		param["userId"] = strconv.FormatInt(user.Id, 10)
	}
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.QQ {
		return web.JsonErrorMsg("QQ登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	param["ref"] = ref
	url := qq.AuthorizeUrl(param)
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
	userId := c.Ctx.FormValue("userId")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByQQ(code, state, userId)
	if err != nil {
		return web.JsonError(err)
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return web.JsonError(codeErr)
	} else {
		return render.BuildLoginSuccess(user, "")
	}
}
