package common

import "bbs-go/simple"

var (
	CaptchaError = simple.NewError(1000, "验证码错误")
)
