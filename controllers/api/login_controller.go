package api

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/utils/github"
)

type LoginController struct {
	Ctx iris.Context
}

// 退出登录
func (this *LoginController) GetSignout() *simple.JsonResult {
	err := services.UserTokenService.Signout(this.Ctx)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonSuccess()
}

// 获取Github授权地址
func (this *LoginController) GetGithub() *simple.JsonResult {
	url := github.OauthConfig.AuthCodeURL(simple.Uuid())
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取Github回调信息获取
func (this *LoginController) GetGithubCallback() *simple.JsonResult {
	code := this.Ctx.FormValue("code")
	githubUser, err := services.GithubUserService.GetGithubUser(code)
	if err != nil {
		logrus.Errorf("Code exchange failed with '%s'", err)
		return simple.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByGithub(githubUser)
	if codeErr != nil {
		return simple.JsonError(codeErr)
	} else { // 直接登录
		token, err := services.UserTokenService.Generate(user.Id)
		if err != nil {
			return simple.JsonErrorMsg(err.Error())
		}
		return simple.NewEmptyRspBuilder().
			Put("token", token).
			Put("user", render.BuildUser(user)).
			JsonResult()
	}
}
