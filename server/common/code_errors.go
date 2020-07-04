package common

import "github.com/mlogclub/simple"

var (
	CaptchaError   = simple.NewError(1000, "验证码错误")
	ForbiddenError = simple.NewError(1001, "已被禁言")
)
