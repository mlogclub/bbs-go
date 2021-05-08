package api

import (
	"bbs-go/controllers/render"
	"bbs-go/pkg/osc"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type OscLoginController struct {
	Ctx iris.Context
}

// GetAuthorize 获取登录授权地址
func (c *OscLoginController) GetAuthorize() *simple.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Osc {
		return simple.JsonErrorMsg("开源中国账号登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	url := osc.AuthCodeURL(map[string]string{"ref": ref})
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// GetCallback 获取回调信息获取
func (c *OscLoginController) GetCallback() *simple.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Osc {
		return simple.JsonErrorMsg("开源中国账号登录/注册已禁用")
	}

	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByOSC(code, state)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return simple.JsonError(codeErr)
	} else {
		return render.BuildLoginSuccess(user, "")
	}
}
