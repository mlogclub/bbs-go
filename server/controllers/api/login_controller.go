package api

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/github"
	"github.com/mlogclub/bbs-go/common/qq"
	"github.com/mlogclub/bbs-go/services"
)

type LoginController struct {
	Ctx iris.Context
}

// 用户名密码登录
func (this *LoginController) PostSignin() *simple.JsonResult {
	var (
		username = this.Ctx.PostValueTrim("username")
		password = this.Ctx.PostValueTrim("password")
		ref      = this.Ctx.FormValue("ref")
	)
	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return services.UserTokenService.GenerateResult(user, ref)
}

// 退出登录
func (this *LoginController) GetSignout() *simple.JsonResult {
	err := services.UserTokenService.Signout(this.Ctx)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 获取Github登录授权地址
func (this *LoginController) GetGithubAuthorize() *simple.JsonResult {
	ref := this.Ctx.FormValue("ref")
	url := github.GetOauthConfig(map[string]string{"ref": ref}).AuthCodeURL(simple.Uuid())
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取Github回调信息获取
func (this *LoginController) GetGithubCallback() *simple.JsonResult {
	code := this.Ctx.FormValue("code")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByGithub(code)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return simple.JsonError(codeErr)
	} else {
		return services.UserTokenService.GenerateResult(user, "")
	}
}

// 获取QQ登录授权地址
func (this *LoginController) GetQqAuthorize() *simple.JsonResult {
	ref := this.Ctx.FormValue("ref")
	url := qq.GetOauthConfig(map[string]string{"ref": ref}).AuthCodeURL(simple.Uuid())
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
