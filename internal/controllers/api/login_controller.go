package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/models"
	"database/sql"

	"github.com/dchest/captcha"
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/pkg/bbsurls"
	captcha2 "bbs-go/internal/pkg/captcha"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/google"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/services"
)

type LoginController struct {
	Ctx iris.Context
}

// 注册
func (c *LoginController) PostSignup() *web.JsonResult {
	var (
		captchaId          = c.Ctx.PostValueTrim("captchaId")
		captchaCode        = c.Ctx.PostValueTrim("captchaCode")
		captchaProtocol, _ = params.GetInt(c.Ctx, "captchaProtocol")
		email              = c.Ctx.PostValueTrim("email")
		username           = c.Ctx.PostValueTrim("username")
		password           = c.Ctx.PostValueTrim("password")
		rePassword         = c.Ctx.PostValueTrim("rePassword")
		nickname           = c.Ctx.PostValueTrim("nickname")
		redirect           = c.Ctx.FormValue("redirect")
	)
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.password_login_disabled"))
	}
	// 根据验证码协议版本校验验证码
	if captchaProtocol == 2 {
		if !captcha2.Verify(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	} else {
		if !captcha.VerifyString(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	}
	user, err := services.UserService.SignUp(username, email, nickname, password, rePassword)
	if err != nil {
		return web.JsonError(err)
	}
	return render.BuildLoginSuccess(c.Ctx, user, redirect)
}

// 用户名密码登录
func (c *LoginController) PostSignin() *web.JsonResult {
	var (
		captchaId          = c.Ctx.PostValueTrim("captchaId")
		captchaCode        = c.Ctx.PostValueTrim("captchaCode")
		captchaProtocol, _ = params.GetInt(c.Ctx, "captchaProtocol")
		username           = c.Ctx.PostValueTrim("username")
		password           = c.Ctx.PostValueTrim("password")
		redirect           = c.Ctx.FormValue("redirect")
	)

	// 根据验证码协议版本校验验证码
	if captchaProtocol == 2 {
		if !captcha2.Verify(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	} else {
		if !captcha.VerifyString(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	}

	user, err := services.UserService.SignIn(username, password)
	if err != nil {
		return web.JsonError(err)
	}

	// 管理员可以突破密码登录的限制，因为后台只能密码登录
	if !user.IsOwnerOrAdmin() {
		if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
			return web.JsonErrorMsg(locales.Get("auth.password_login_disabled"))
		}
	}
	return render.BuildLoginSuccess(c.Ctx, user, redirect)
}

// 请求找回密码邮件
func (c *LoginController) PostSend_reset_password_email() *web.JsonResult {
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.password_login_disabled"))
	}
	var (
		captchaId          = c.Ctx.PostValueTrim("captchaId")
		captchaCode        = c.Ctx.PostValueTrim("captchaCode")
		captchaProtocol, _ = params.GetInt(c.Ctx, "captchaProtocol")
		email              = c.Ctx.PostValueTrim("email")
	)

	// 根据验证码协议版本校验验证码
	if captchaProtocol == 2 {
		if !captcha2.Verify(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	} else {
		if !captcha.VerifyString(captchaId, captchaCode) {
			return web.JsonError(errs.CaptchaError())
		}
	}

	if err := services.UserService.SendResetPasswordEmail(email); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 重置密码
func (c *LoginController) PostReset_password() *web.JsonResult {
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.password_login_disabled"))
	}
	var (
		token      = c.Ctx.PostValueTrim("token")
		password   = c.Ctx.PostValueTrim("password")
		rePassword = c.Ctx.PostValueTrim("rePassword")
	)
	if err := services.UserService.ResetPasswordByToken(token, password, rePassword); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 退出登录
func (c *LoginController) GetSignout() *web.JsonResult {
	err := services.UserTokenService.Signout(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 请求登录短信验证码
func (c *LoginController) PostLogin_sms_code() *web.JsonResult {
	var (
		phone, _       = params.Get(c.Ctx, "phone")
		captchaId, _   = params.Get(c.Ctx, "captchaId")
		captchaCode, _ = params.Get(c.Ctx, "captchaCode")
	)

	if strs.IsBlank(phone) {
		return web.JsonErrorMsg(locales.Get("auth.phone_required"))
	}

	if !captcha2.Verify(captchaId, captchaCode) {
		return web.JsonError(errs.CaptchaError())
	}

	if !services.SysConfigService.GetLoginConfig().SmsLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.sms_login_disabled"))
	}

	smsId, err := services.SmsCodeService.SendSms(phone)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(iris.Map{
		"smsId": smsId,
	})
}

// 短信登录
func (c *LoginController) PostLogin_sms() *web.JsonResult {
	var (
		smsId, _   = params.Get(c.Ctx, "smsId")
		smsCode, _ = params.Get(c.Ctx, "smsCode")
		redirect   = c.Ctx.FormValue("redirect")
	)

	if !services.SysConfigService.GetLoginConfig().SmsLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.sms_login_disabled"))
	}

	phone, err := services.SmsCodeService.Verify(smsId, smsCode)
	if err != nil {
		return web.JsonError(err)
	}

	user := services.UserService.GetByPhone(phone)
	if user == nil {
		user = &models.User{
			Phone: sql.NullString{
				String: phone,
				Valid:  true,
			},
			Nickname:   "User" + common.StrRight(phone, 4),
			CreateTime: dates.NowTimestamp(),
			UpdateTime: dates.NowTimestamp(),
		}

		if err := services.UserService.Create(user); err != nil {
			return web.JsonError(err)
		}
	}

	return render.BuildLoginSuccess(c.Ctx, user, redirect)
}

func (c *LoginController) GetWx_login_config() *web.JsonResult {
	redirect, _ := params.Get(c.Ctx, "redirect")
	bind, _ := params.GetBool(c.Ctx, "bind")
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.WeixinLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled"))
	}

	cache.WxLoginStateCache.Put(state, &cache.WxLoginStateData{
		Redirect: redirect,
		Bind:     bind,
	})

	redirectURI := bbsurls.AbsUrl("/user/signin/callback/weixin")
	if bind {
		redirectURI = bbsurls.AbsUrl("/user/signin/callback/weixin_bind")
	}

	return web.JsonData(iris.Map{
		"appid":        loginConfig.WeixinLogin.AppId,
		"scope":        "snsapi_login",
		"redirect_uri": redirectURI,
		"state":        state,
	})
}

func (c *LoginController) PostWx_login_submit() *web.JsonResult {
	code, _ := params.Get(c.Ctx, "code")
	state, _ := params.Get(c.Ctx, "state")

	data := cache.WxLoginStateCache.Get(state)
	if data == nil {
		return web.JsonErrorMsg(locales.Get("auth.login_data_error"))
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled"))
	}

	user, err := services.ThirdUserService.LoginWeixin(code, state)
	if err != nil {
		return web.JsonError(err)
	}

	return render.BuildLoginSuccess(c.Ctx, user, data.Redirect)
}

func (c *LoginController) PostWx_bind() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	code, _ := params.Get(c.Ctx, "code")
	state, _ := params.Get(c.Ctx, "state")

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled"))
	}

	if err := services.ThirdUserService.BindWeixin(user.Id, code, state); err != nil {
		return web.JsonError(err)
	}

	return web.JsonSuccess()
}

func (c *LoginController) PostWx_unbind() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		return web.JsonErrorMsg(locales.Get("auth.weixin_login_disabled"))
	}

	services.ThirdUserService.UnbindWeixin(user.Id)

	return web.JsonSuccess()
}

func (c *LoginController) GetGoogle_login_config() *web.JsonResult {
	redirect, _ := params.Get(c.Ctx, "redirect")
	bind, _ := params.GetBool(c.Ctx, "bind")
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		return web.JsonErrorMsg("Google登录未启用")
	}

	cache.GoogleLoginStateCache.Put(state, &cache.GoogleLoginStateData{
		Redirect: redirect,
		Bind:     bind,
	})

	redirectURI := bbsurls.AbsUrl(google.CallbackPathLogin)
	if bind {
		redirectURI = bbsurls.AbsUrl(google.CallbackPathBind)
	}

	oauth := google.NewGoogleOAuth(loginConfig.GoogleLogin.ClientId, loginConfig.GoogleLogin.ClientSecret, redirectURI)
	authURL := oauth.GetAuthURL(state)

	return web.JsonData(iris.Map{
		"clientId":    loginConfig.GoogleLogin.ClientId,
		"authUrl":     authURL,
		"redirectUri": redirectURI, // 用于调试，显示实际使用的 redirect URI
		"state":       state,
		"redirect":    redirect,
	})
}

func (c *LoginController) PostGoogle_login_submit() *web.JsonResult {
	code, _ := params.Get(c.Ctx, "code")
	state, _ := params.Get(c.Ctx, "state")

	if strs.IsBlank(state) {
		return web.JsonErrorMsg("state参数缺失")
	}

	data := cache.GoogleLoginStateCache.Get(state)
	if data == nil {
		return web.JsonErrorMsg("登录数据错误或已过期，请重新登录")
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		return web.JsonErrorMsg("Google登录未启用")
	}

	user, err := services.ThirdUserService.LoginGoogle(code, state)
	if err != nil {
		return web.JsonError(err)
	}

	return render.BuildLoginSuccess(c.Ctx, user, data.Redirect)
}

func (c *LoginController) PostGoogle_bind() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	code, _ := params.Get(c.Ctx, "code")
	state, _ := params.Get(c.Ctx, "state")

	if strs.IsBlank(state) {
		return web.JsonErrorMsg("state参数缺失")
	}

	// 验证 state 是否存在（防止 CSRF 攻击）
	data := cache.GoogleLoginStateCache.Get(state)
	if data == nil {
		return web.JsonErrorMsg("绑定数据错误或已过期，请重新绑定")
	}

	// 验证是否为绑定流程
	if !data.Bind {
		return web.JsonErrorMsg("无效的绑定请求")
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		return web.JsonErrorMsg("Google登录未启用")
	}

	if err := services.ThirdUserService.BindGoogle(user.Id, code, state); err != nil {
		return web.JsonError(err)
	}

	return web.JsonSuccess()
}

func (c *LoginController) PostGoogle_one_tap() *web.JsonResult {
	credential, _ := params.Get(c.Ctx, "credential")
	if credential == "" {
		return web.JsonErrorMsg("credential参数缺失")
	}

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		return web.JsonErrorMsg("Google登录未启用")
	}

	user, err := services.ThirdUserService.LoginGoogleOneTap(credential)
	if err != nil {
		return web.JsonError(err)
	}

	return render.BuildLoginSuccess(c.Ctx, user, "")
}

func (c *LoginController) PostGoogle_unbind() *web.JsonResult {
	user, err := common.CheckLogin(c.Ctx)
	if err != nil {
		return web.JsonError(err)
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		return web.JsonErrorMsg("Google登录未启用")
	}

	services.ThirdUserService.UnbindGoogle(user.Id)

	return web.JsonSuccess()
}
