package api

import (
	"bbs-go/internal/controllers/render"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/validate"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/services"
)

type ForgetPasswordController struct {
	Ctx iris.Context
}

// 定义用于JSON请求的结构体
type SendEmailRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token      string `json:"token"`
	Password   string `json:"password"`
	RePassword string `json:"rePassword"`
}

// 发送密码重置邮件
func (c *ForgetPasswordController) PostSendEmail() *web.JsonResult {
	var email string

	// 检查请求是否为JSON格式 (注意包含charset)
	contentType := c.Ctx.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// JSON请求格式
		var req SendEmailRequest
		if err := c.Ctx.ReadJSON(&req); err != nil {
			return web.JsonError(err)
		}
		email = strings.TrimSpace(req.Email)
	} else {
		// 表单请求格式
		email = strings.TrimSpace(c.Ctx.PostValueTrim("email"))
	}

	if strs.IsBlank(email) {
		return web.JsonErrorMsg(locales.Get("errors.email_empty"))
	}
	if err := validate.IsEmail(email); err != nil {
		// 将验证错误消息替换为多语言化的消息
		return web.JsonErrorMsg(locales.Get("errors.email_invalid"))
	}

	err := services.UserService.SendPasswordResetEmail(email)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}

// 重置密码
func (c *ForgetPasswordController) PostResetPassword() *web.JsonResult {
	var (
		token      string
		password   string
		rePassword string
	)

	// 检查请求是否为JSON格式 (注意包含charset)
	contentType := c.Ctx.GetHeader("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// JSON请求格式
		var req ResetPasswordRequest
		if err := c.Ctx.ReadJSON(&req); err != nil {
			return web.JsonError(err)
		}
		token = req.Token
		password = req.Password
		rePassword = req.RePassword
	} else {
		// 表单请求格式
		token = c.Ctx.PostValueTrim("token")
		password = c.Ctx.PostValueTrim("password")
		rePassword = c.Ctx.PostValueTrim("rePassword")
	}

	if strs.IsBlank(token) {
		return web.JsonErrorMsg(locales.Get("errors.reset_link_invalid"))
	}
	if strs.IsBlank(password) {
		return web.JsonErrorMsg(locales.Get("errors.password_empty"))
	}
	if strs.IsBlank(rePassword) {
		return web.JsonErrorMsg(locales.Get("errors.confirm_password_empty"))
	}

	err := services.UserService.ResetPassword(token, password, rePassword)
	if err != nil {
		return web.JsonError(err)
	}

	// 重置成功后，可能需要返回登录信息
	user, err := services.UserService.GetUserByPasswordResetToken(token)
	if err != nil {
		return web.JsonError(err)
	}

	// 标记重置链接已使用
	err = services.UserService.MarkPasswordResetTokenUsed(token)
	if err != nil {
		return web.JsonError(err)
	}

	return render.BuildLoginSuccess(c.Ctx, user, "")
}