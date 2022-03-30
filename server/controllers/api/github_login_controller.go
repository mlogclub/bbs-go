package api

import (
	"bbs-go/controllers/render"
	"bbs-go/pkg/github"
	"bbs-go/services"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"
)

type GithubLoginController struct {
	Ctx iris.Context
}

// 获取Github登录授权地址
func (c *GithubLoginController) GetAuthorize() *web.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Github {
		return web.JsonErrorMsg("Github登录/注册已禁用")
	}

	ref := c.Ctx.FormValue("ref")
	url := github.AuthCodeURL(map[string]string{"ref": ref})
	return web.NewEmptyRspBuilder().Put("url", url).JsonResult()
}

// 获取Github回调信息获取
func (c *GithubLoginController) GetCallback() *web.JsonResult {
	loginMethod := services.SysConfigService.GetLoginMethod()
	if !loginMethod.Github {
		return web.JsonErrorMsg("Github登录/注册已禁用")
	}

	code := c.Ctx.FormValue("code")
	state := c.Ctx.FormValue("state")

	thirdAccount, err := services.ThirdAccountService.GetOrCreateByGithub(code, state)
	if err != nil {
		return web.JsonErrorMsg(err.Error())
	}

	user, codeErr := services.UserService.SignInByThirdAccount(thirdAccount)
	if codeErr != nil {
		return web.JsonError(codeErr)
	} else {
		return render.BuildLoginSuccess(user, "")
	}
}
