package api

import (
	"bbs-go/controllers/render"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/mvc"

	"bbs-go/pkg/common"
	"bbs-go/services"
)

type LoginController struct {
	Ctx iris.Context
}

// 注册
func (c *LoginController) PostSignup() *mvc.JsonResult {
	var (
		captchaId   = c.Ctx.PostValueTrim("captchaId")
		captchaCode = c.Ctx.PostValueTrim("captchaCode")
		email       = c.Ctx.PostValueTrim("email")
		username    = c.Ctx.PostValueTrim("username")
		password    = c.Ctx.PostValueTrim("password")
		rePassword  = c.Ctx.PostValueTrim("rePassword")
		nickname    = c.Ctx.PostValueTrim("nickname")
		ref         = c.Ctx.FormValue("ref")
	)
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Password {
		return mvc.JsonErrorMsg("账号密码登录/注册已禁用")
	}
	if !captcha.VerifyString(captchaId, captchaCode) {
		return mvc.JsonError(common.CaptchaError)
	}
	user, err := services.UserService.SignUp(username, email, nickname, password, rePassword)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return render.BuildLoginSuccess(user, ref)
}

// 用户名密码登录
func (c *LoginController) PostSignin() *mvc.JsonResult {
	var (
		captchaId   = c.Ctx.PostValueTrim("captchaId")
		captchaCode = c.Ctx.PostValueTrim("captchaCode")
		username    = c.Ctx.PostValueTrim("username")
		password    = c.Ctx.PostValueTrim("password")
		ref         = c.Ctx.FormValue("ref")
	)
	// loginMethod := services.SysConfigService.GetLoginMethod()
	// if !loginMethod.Password {
	// 	return mvc.JsonErrorMsg("账号密码登录/注册已禁用")
	// }
	if !captcha.VerifyString(captchaId, captchaCode) {
		return mvc.JsonError(common.CaptchaError)
	}
	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return render.BuildLoginSuccess(user, ref)
}

// 退出登录
func (c *LoginController) GetSignout() *mvc.JsonResult {
	err := services.UserTokenService.Signout(c.Ctx)
	if err != nil {
		return mvc.JsonErrorMsg(err.Error())
	}
	return mvc.JsonSuccess()
}
