package api

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/req"
	"bbs-go/internal/services/heatpoints"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"bbs-go/internal/pkg/ginx"

	"github.com/dchest/captcha"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"

	"bbs-go/internal/pkg/bbsurls"
	captcha2 "bbs-go/internal/pkg/captcha"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/pkg/github"
	"bbs-go/internal/pkg/google"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/params"
	"bbs-go/internal/services"
)

// 注册
// 用户名密码登录
// 请求找回密码邮件
// 重置密码
// 退出登录
// 请求登录短信验证码
// 短信登录
func LoginSignup(ctx *gin.Context) {
	var req req.LoginSignupReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	req.CaptchaId = strings.TrimSpace(req.CaptchaId)
	req.CaptchaCode = strings.TrimSpace(req.CaptchaCode)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.RePassword = strings.TrimSpace(req.RePassword)
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Redirect == "" {
		req.Redirect = ctx.Query("redirect")
	}
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.password_login_disabled")))
		return
	}
	// 根据验证码协议版本校验验证码
	if req.CaptchaProtocol == 2 {
		if !captcha2.Verify(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	} else {
		if !captcha.VerifyString(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	}
	user, err := services.UserService.SignUp(req.Username, req.Email, req.Nickname, req.Password, req.RePassword)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	// 新用户空投热度点
	if err := heatpoints.Airdrop.AirdropSingleUser(user.Id); err != nil {
		slog.Warn("新用户空投失败", "userId", user.Id, "error", err)
	}
	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, req.Redirect))

}

func LoginSignin(ctx *gin.Context) {
	var req req.LoginSigninReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	req.CaptchaId = strings.TrimSpace(req.CaptchaId)
	req.CaptchaCode = strings.TrimSpace(req.CaptchaCode)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Redirect == "" {
		req.Redirect = ctx.Query("redirect")
	}

	// 根据验证码协议版本校验验证码
	if req.CaptchaProtocol == 2 {
		if !captcha2.Verify(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	} else {
		if !captcha.VerifyString(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	}

	user, err := services.UserService.SignIn(req.Username, req.Password)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	// 站长可以突破密码登录的限制，因为后台只能密码登录
	if !user.IsOwner() {
		if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
			ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.password_login_disabled")))
			return
		}
	}
	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, req.Redirect))

}

func LoginSendResetPasswordEmail(ctx *gin.Context) {
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.password_login_disabled")))
		return
	}
	var req req.LoginResetEmailReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	req.CaptchaId = strings.TrimSpace(req.CaptchaId)
	req.CaptchaCode = strings.TrimSpace(req.CaptchaCode)
	req.Email = strings.TrimSpace(req.Email)

	// 根据验证码协议版本校验验证码
	if req.CaptchaProtocol == 2 {
		if !captcha2.Verify(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	} else {
		if !captcha.VerifyString(req.CaptchaId, req.CaptchaCode) {
			ginx.WriteJSON(ctx, errs.CaptchaError())
			return
		}
	}

	if err := services.UserService.SendResetPasswordEmail(req.Email); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func LoginResetPassword(ctx *gin.Context) {
	if !services.SysConfigService.GetLoginConfig().PasswordLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.password_login_disabled")))
		return
	}
	var req req.LoginResetPasswordReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	req.Token = strings.TrimSpace(req.Token)
	req.Password = strings.TrimSpace(req.Password)
	req.RePassword = strings.TrimSpace(req.RePassword)
	if err := services.UserService.ResetPasswordByToken(req.Token, req.Password, req.RePassword); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func LoginSignout(ctx *gin.Context) {
	err := services.UserTokenService.Signout(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, nil)

}

func LoginLoginSmsCode(ctx *gin.Context) {
	var req req.LoginSmsCodeReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if strs.IsBlank(req.Phone) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.phone_required")))
		return
	}

	if !captcha2.Verify(req.CaptchaId, req.CaptchaCode) {
		ginx.WriteJSON(ctx, errs.CaptchaError())
		return
	}

	if !services.SysConfigService.GetLoginConfig().SmsLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.sms_login_disabled")))
		return
	}

	smsId, err := services.SmsCodeService.SendSms(req.Phone)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	ginx.WriteJSON(ctx, map[string]interface{}{
		"smsId": smsId,
	})

}

func LoginLoginSms(ctx *gin.Context) {
	var req req.LoginSmsReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	if req.Redirect == "" {
		req.Redirect = ctx.Query("redirect")
	}

	if !services.SysConfigService.GetLoginConfig().SmsLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.sms_login_disabled")))
		return
	}

	phone, err := services.SmsCodeService.Verify(req.SmsId, req.SmsCode)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
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
			ginx.WriteJSON(ctx, err)
			return
		}
	}

	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, req.Redirect))

}

func LoginWxLoginConfig(ctx *gin.Context) {
	var req req.OAuthConfigReq
	if err := ginx.BindQuery(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.WeixinLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.weixin_login_disabled")))
		return
	}

	cache.WxLoginStateCache.Put(state, &cache.WxLoginStateData{
		Redirect: req.Redirect,
		Bind:     req.Bind,
	})

	redirectURI := bbsurls.AbsUrl("/user/signin/callback/weixin")
	if req.Bind {
		redirectURI = bbsurls.AbsUrl("/user/signin/callback/weixin_bind")
	}

	ginx.WriteJSON(ctx, map[string]interface{}{
		"appid":        loginConfig.WeixinLogin.AppId,
		"scope":        "snsapi_login",
		"redirect_uri": redirectURI,
		"state":        state,
	})

}

func LoginWxLoginSubmit(ctx *gin.Context) {
	var req req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	data := cache.WxLoginStateCache.Get(req.State)
	if data == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.login_data_error")))
		return
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.weixin_login_disabled")))
		return
	}

	user, err := services.ThirdUserService.LoginWeixin(req.Code, req.State)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, data.Redirect))

}

func LoginWxBind(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	var req req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.weixin_login_disabled")))
		return
	}

	if err := services.ThirdUserService.BindWeixin(user.Id, req.Code, req.State); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, nil)

}

func LoginWxUnbind(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if !services.SysConfigService.GetLoginConfig().WeixinLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.weixin_login_disabled")))
		return
	}

	services.ThirdUserService.UnbindWeixin(user.Id)

	ginx.WriteJSON(ctx, nil)

}

func LoginGoogleLoginConfig(ctx *gin.Context) {
	var req req.OAuthConfigReq
	if err := ginx.BindQuery(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.google_login_disabled")))
		return
	}

	cache.GoogleLoginStateCache.Put(state, &cache.GoogleLoginStateData{
		Redirect: req.Redirect,
		Bind:     req.Bind,
	})

	redirectURI := bbsurls.AbsUrl(google.CallbackPathLogin)
	if req.Bind {
		redirectURI = bbsurls.AbsUrl(google.CallbackPathBind)
	}

	oauth := google.NewGoogleOAuth(loginConfig.GoogleLogin.ClientId, loginConfig.GoogleLogin.ClientSecret, redirectURI)
	authURL := oauth.GetAuthURL(state)

	ginx.WriteJSON(ctx, map[string]interface{}{
		"clientId":    loginConfig.GoogleLogin.ClientId,
		"authUrl":     authURL,
		"redirectUri": redirectURI, // 用于调试，显示实际使用的 redirect URI
		"state":       state,
		"redirect":    req.Redirect,
	})

}

func LoginGoogleLoginSubmit(ctx *gin.Context) {
	var req req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if strs.IsBlank(req.State) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.state_required")))
		return
	}

	data := cache.GoogleLoginStateCache.Get(req.State)
	if data == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.login_data_error")))
		return
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.google_login_disabled")))
		return
	}

	user, err := services.ThirdUserService.LoginGoogle(req.Code, req.State)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, data.Redirect))

}

func LoginGoogleBind(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	var req req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if strs.IsBlank(req.State) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.state_required")))
		return
	}

	// 验证 state 是否存在（防止 CSRF 攻击）
	data := cache.GoogleLoginStateCache.Get(req.State)
	if data == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.bind_data_error")))
		return
	}

	// 验证是否为绑定流程
	if !data.Bind {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.invalid_bind_request")))
		return
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.google_login_disabled")))
		return
	}

	if err := services.ThirdUserService.BindGoogle(user.Id, req.Code, req.State); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, nil)

}

func LoginGoogleOneTap(ctx *gin.Context) {
	credential := params.FormValue(ctx, "credential")
	if credential == "" {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.credential_required")))
		return
	}

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GoogleLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.google_login_disabled")))
		return
	}

	user, err := services.ThirdUserService.LoginGoogleOneTap(credential)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, ""))

}

func LoginGoogleUnbind(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if !services.SysConfigService.GetLoginConfig().GoogleLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.google_login_disabled")))
		return
	}

	services.ThirdUserService.UnbindGoogle(user.Id)

	ginx.WriteJSON(ctx, nil)

}

func LoginGithubLoginConfig(ctx *gin.Context) {
	var req req.OAuthConfigReq
	if err := ginx.BindQuery(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}
	state := strs.UUID()

	loginConfig := services.SysConfigService.GetLoginConfig()
	if !loginConfig.GithubLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.github_login_disabled")))
		return
	}

	// 绑定流程默认跳转到账号设置页，登录流程则使用传入的 redirect
	if req.Bind && strs.IsBlank(req.Redirect) {
		req.Redirect = "/user/profile/account"
	}

	cache.GithubLoginStateCache.Put(state, &cache.GithubLoginStateData{
		Redirect: req.Redirect,
		Bind:     req.Bind,
	})

	// GitHub OAuth 只支持配置一个 callback URL，这里统一使用登录回调路径
	redirectURI := bbsurls.AbsUrl(github.AuthorizationCallbackURL)

	oauth := github.NewGithubOAuth(loginConfig.GithubLogin.ClientId, loginConfig.GithubLogin.ClientSecret, redirectURI)
	authURL := oauth.GetAuthURL(state)

	ginx.WriteJSON(ctx, map[string]interface{}{
		"clientId":    loginConfig.GithubLogin.ClientId,
		"authUrl":     authURL,
		"redirectUri": redirectURI,
		"state":       state,
		"redirect":    req.Redirect,
	})

}

func LoginGithubLoginSubmit(ctx *gin.Context) {
	var req req.OAuthCodeStateReq
	if err := ginx.Bind(ctx, &req); err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if strs.IsBlank(req.State) {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.state_required")))
		return
	}

	data := cache.GithubLoginStateCache.Get(req.State)
	if data == nil {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.login_data_error")))
		return
	}

	if !services.SysConfigService.GetLoginConfig().GithubLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.github_login_disabled")))
		return
	}

	// 根据 state 中的标记区分「登录」和「绑定」场景
	if data.Bind {
		// 绑定流程要求用户已登录
		user, err := common.CheckLogin(ctx)
		if err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}
		if err := services.ThirdUserService.BindGithub(user.Id, req.Code, req.State); err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}
		// 绑定成功后保持当前登录态不变，按 redirect 跳转（默认为账号设置页）
		ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, data.Redirect))
		return
	} else {
		// 普通登录流程：使用 GitHub 账号登录 / 注册
		user, err := services.ThirdUserService.LoginGithub(req.Code, req.State)
		if err != nil {
			ginx.WriteJSON(ctx, err)
			return
		}

		ginx.WriteJSON(ctx, render.BuildLoginSuccess(ctx, user, data.Redirect))

		return
	}

}

func LoginGithubUnbind(ctx *gin.Context) {
	user, err := common.CheckLogin(ctx)
	if err != nil {
		ginx.WriteJSON(ctx, err)
		return
	}

	if !services.SysConfigService.GetLoginConfig().GithubLogin.Enabled {
		ginx.WriteJSON(ctx, ginx.ErrorMessage(locales.Get("auth.github_login_disabled")))
		return
	}

	services.ThirdUserService.UnbindGithub(user.Id)
	ginx.WriteJSON(ctx, nil)

}
