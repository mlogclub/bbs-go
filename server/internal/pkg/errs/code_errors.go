package errs

import (
	"bbs-go/internal/pkg/locales"

	"github.com/mlogclub/simple/web"
)

const (
	CodeNotLogin            = 1
	CodeCaptchaError        = 1000
	CodeForbiddenError      = 1001
	CodeUserDisabled        = 1002
	CodeInObservationPeriod = 1003
	CodeEmailNotVerified    = 1004
)

// NewError 创建错误
func NewError(code int) *web.CodeError {
	var message string
	switch code {
	case CodeNotLogin:
		message = locales.Get("errors.not_login")
	case CodeCaptchaError:
		message = locales.Get("errors.captcha_error")
	case CodeForbiddenError:
		message = locales.Get("errors.forbidden")
	case CodeUserDisabled:
		message = locales.Get("errors.user_disabled")
	case CodeEmailNotVerified:
		message = locales.Get("errors.email_not_verified")
	default:
		message = "Unknown error"
	}
	return web.NewError(code, message)
}

// 预定义的错误创建函数
var (
	NotLogin         = func() *web.CodeError { return NewError(CodeNotLogin) }
	CaptchaError     = func() *web.CodeError { return NewError(CodeCaptchaError) }
	ForbiddenError   = func() *web.CodeError { return NewError(CodeForbiddenError) }
	UserDisabled     = func() *web.CodeError { return NewError(CodeUserDisabled) }
	EmailNotVerified = func() *web.CodeError { return NewError(CodeEmailNotVerified) }
)
