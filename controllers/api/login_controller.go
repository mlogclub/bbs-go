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

func (this *LoginController) GetGithub() *simple.JsonResult {
	url := github.OauthConfig.AuthCodeURL(simple.Uuid())
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

func (this *LoginController) GetGithubCallback() *simple.JsonResult {
	code := this.Ctx.FormValue("code")
	githubUser, err := services.GithubUserService.GetGithubUser(code)
	if err != nil {
		logrus.Errorf("Code exchange failed with '%s'", err)
		return simple.ErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByGithub(githubUser)
	if codeErr != nil {
		return simple.Error(codeErr)
	} else { // 直接登录
		token, err := services.UserTokenService.Generate(user.Id)
		if err != nil {
			return simple.ErrorMsg(err.Error())
		}
		return simple.NewEmptyRspBuilder().
			Put("token", token).
			Put("user", render.BuildUser(user)).
			JsonResult()
	}
}
