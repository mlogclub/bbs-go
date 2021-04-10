package api

import (
	"bbs-go/controllers/render"
	"bbs-go/package/github"
	"bbs-go/services"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple"
)

type GithubLoginController struct {
	Ctx iris.Context
}

// 获取Github登录授权地址
func (c *GithubLoginController) GetAuthorize() *simple.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Github {
		return simple.JsonErrorMsg("Github登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	url := github.AuthCodeURL(map[string]string{"ref": ref})
	return simple.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取Github回调信息获取
func (c *GithubLoginController) GetCallback() *simple.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Github {
		return simple.JsonErrorMsg("Github登录/注册已禁用")
	}

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
		return render.BuildLoginSuccess(user, "")
	}
}
