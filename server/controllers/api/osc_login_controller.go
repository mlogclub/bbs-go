package api

import (
	"bbs-go/controllers/render"
	"bbs-go/pkg/osc"
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"
)

type OscLoginController struct {
	Ctx iris.Context
}

// GetAuthorize 获取登录授权地址
func (c *OscLoginController) GetAuthorize() *mvc.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Osc {
		return mvc.JsonErrorMsg("开源中国账号登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	url := osc.AuthCodeURL(map[string]string{"ref": ref})
	return mvc.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// GetCallback 获取回调信息获取
func (c *OscLoginController) GetCallback() *mvc.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Osc {
		return mvc.JsonErrorMsg("开源中国账号登录/注册已禁用")
	}

	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByOSC(code, state)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return mvc.JsonError(codeErr)
	} else {
		return render.BuildLoginSuccess(user, "")
	}
}
