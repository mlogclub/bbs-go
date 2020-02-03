package api

import (
	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"

	"bbs-go/common"
	"bbs-go/common/github"
	"bbs-go/common/qq"
	"bbs-go/controllers/render"
	"bbs-go/model"
	"bbs-go/services"
)

type LoginController struct {
	Ctx iris.Context
}

// 注册
func (c *LoginController) PostSignup() *simple.JsonResult {
	var (
		captchaId   = c.Ctx.PostValueTrim("captchaId")
		captchaCode = c.Ctx.PostValueTrim("captchaCode")
		username    = c.Ctx.PostValueTrim("username")
		password    = c.Ctx.PostValueTrim("password")
		rePassword  = c.Ctx.PostValueTrim("rePassword")
		nickname    = c.Ctx.PostValueTrim("nickname")
		ref         = c.Ctx.FormValue("ref")
	)
	if !captcha.VerifyString(captchaId, captchaCode) {
		return simple.JsonError(common.CaptchaError)
	}
	user, err := services.UserService.SignUp(username, "", nickname, password, rePassword)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return c.GenerateLoginResult(user, ref)
}

// 用户名密码登录
func (c *LoginController) PostSignin() *simple.JsonResult {
	var (
		captchaId   = c.Ctx.PostValueTrim("captchaId")
		captchaCode = c.Ctx.PostValueTrim("captchaCode")
		username    = c.Ctx.PostValueTrim("username")
		password    = c.Ctx.PostValueTrim("password")
		ref         = c.Ctx.FormValue("ref")
	)
	if !captcha.VerifyString(captchaId, captchaCode) {
		return simple.JsonError(common.CaptchaError)
	}
	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return c.GenerateLoginResult(user, ref)
}

// 退出登录
func (c *LoginController) GetSignout() *simple.JsonResult {
	err := services.UserTokenService.Signout(c.Ctx)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 获取Github登录授权地址
func (c *LoginController) GetGithubAuthorize() *simple.JsonResult {
	ref := c.Ctx.FormValue("ref")
	url := github.AuthCodeURL(map[string]string{"ref": ref})
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取Github回调信息获取
func (c *LoginController) GetGithubCallback() *simple.JsonResult {
	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByGithub(code, state)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return simple.JsonError(codeErr)
	} else {
		return c.GenerateLoginResult(user, "")
	}
}

// 获取QQ登录授权地址
func (c *LoginController) GetQqAuthorize() *simple.JsonResult {
	ref := c.Ctx.FormValue("ref")
	url := qq.AuthorizeUrl(map[string]string{"ref": ref})
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取QQ回调信息获取
func (c *LoginController) GetQqCallback() *simple.JsonResult {
	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByQQ(code, state)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return simple.JsonError(codeErr)
	} else {
		return c.GenerateLoginResult(user, "")
	}
}

// user: login user, ref: 登录来源地址，需要控制登录成功之后跳转到该地址
func (c *LoginController) GenerateLoginResult(user *model.User, ref string) *simple.JsonResult {
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", render.BuildUser(user)).
		Put("ref", ref).JsonResult()
}
