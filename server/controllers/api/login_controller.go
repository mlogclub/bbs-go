package api

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/common/github"
	"github.com/mlogclub/bbs-go/common/qq"
	"github.com/mlogclub/bbs-go/controllers/render"
	"github.com/mlogclub/bbs-go/model"
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
	return this.GenerateResult(user, ref)
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
		return this.GenerateResult(user, "")
	}
}

// user: login user, ref: 登录来源地址，需要控制登录成功之后跳转到该地址
func (this *LoginController) GenerateResult(user *model.User, ref string) *simple.JsonResult {
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", render.BuildUser(user)).
		Put("ref", ref).JsonResult()
}

// 获取QQ登录授权地址
func (this *LoginController) GetQqAuthorize() *simple.JsonResult {
	ref := this.Ctx.FormValue("ref")
	url := qq.GetOauthConfig(map[string]string{"ref": ref}).AuthCodeURL(simple.Uuid())
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}
