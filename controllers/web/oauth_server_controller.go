package web

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/mlogclub/mlog/controllers/render"
	"github.com/mlogclub/mlog/services"
	"github.com/mlogclub/mlog/services/oauth"
	"github.com/mlogclub/mlog/utils/session"
	"net/http"

	"github.com/kataras/iris/context"
	"github.com/mlogclub/simple"
)

type OauthServerController struct {
	Ctx context.Context
}

// 请求授权，返回code
func (c *OauthServerController) AnyAuthorize() {
	err := oauth.ServerInstance.Srv.HandleAuthorizeRequest(c.Ctx.ResponseWriter(), c.Ctx.Request())
	if err != nil {
		http.Error(c.Ctx.ResponseWriter(), err.Error(), http.StatusBadRequest)
	}
}

// 通过code换取令牌
func (c *OauthServerController) AnyToken() {
	err := oauth.ServerInstance.Srv.HandleTokenRequest(c.Ctx.ResponseWriter(), c.Ctx.Request())
	if err != nil {
		http.Error(c.Ctx.ResponseWriter(), err.Error(), http.StatusInternalServerError)
	}
}

// 获取用户信息
func (c *OauthServerController) AnyUserinfo() *simple.JsonResult {
	user := oauth.GetUserInfoByRequest(c.Ctx.Request())
	if user == nil {
		return simple.ErrorMsg("User not found")
	}
	return simple.JsonData(user)
}

// 登录页面
func (c *OauthServerController) GetLogin() mvc.View {
	var (
		clientId     = c.Ctx.FormValue("client_id")
		redirectUri  = c.Ctx.FormValue("redirect_uri")
		responseType = c.Ctx.FormValue("response_type")
		state        = c.Ctx.FormValue("state")
	)

	return mvc.View{
		Name: "oauth/login.html",
		Data: iris.Map{
			"client_id":     clientId,
			"redirect_uri":  redirectUri,
			"response_type": responseType,
			"state":         state,
		},
	}
}

// 提交登录
func (c *OauthServerController) PostLogin() {
	var (
		clientId     = c.Ctx.FormValue("client_id")
		redirectUri  = c.Ctx.FormValue("redirect_uri")
		responseType = c.Ctx.FormValue("response_type")
		state        = c.Ctx.FormValue("state")
		username     = c.Ctx.FormValue("username")
		password     = c.Ctx.FormValue("password")
	)
	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		render.View(c.Ctx, "oauth/login.html", iris.Map{
			"client_id":     clientId,
			"redirect_uri":  redirectUri,
			"response_type": responseType,
			"state":         state,
			"ErrMsg":        err.Error(),
			"Username":      username,
			"Password":      password,
		})
		return
	}
	session.SetCurrentUser(c.Ctx, user.Id)

	// TODO 登录成功之后的跳转地址先写死，后面将它配置到数据库中
	c.Ctx.Redirect("/oauth/client", iris.StatusSeeOther)
}
