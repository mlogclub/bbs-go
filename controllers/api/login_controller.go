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

// 获取GithubUser
func (this *LoginController) GetGithubUserBy(githubUserId int64) *simple.JsonResult {
	githubUser := services.GithubUserService.Get(githubUserId)
	if githubUser == nil {
		return simple.JsonErrorMsg("数据不存在")
	}
	return simple.JsonData(githubUser)
}

// Github绑定
func (this *LoginController) PostGithubBind() *simple.JsonResult {
	bindType := this.Ctx.PostValueTrim("bindType")
	githubId, err := this.Ctx.PostValueInt64("githubId")
	username := this.Ctx.PostValueTrim("username")
	email := this.Ctx.PostValueTrim("email")
	password := this.Ctx.PostValueTrim("password")
	rePassword := this.Ctx.PostValueTrim("rePassword")
	nickname := this.Ctx.PostValueTrim("nickname")

	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	// 执行绑定
	user, err := services.UserService.Bind(githubId, bindType, username, email, password, rePassword, nickname)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	// 绑定成功，执行登录
	token, err := services.UserTokenService.Generate(user.Id)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.NewEmptyRspBuilder().
		Put("token", token).
		Put("user", render.BuildUser(user)).
		JsonResult()
}
