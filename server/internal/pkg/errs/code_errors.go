package errs

import (
	"github.com/mlogclub/simple/web"
)

var (
	NotLogin            = web.NewError(1, "请先登录")
	CaptchaError        = web.NewError(1000, "验证码错误")
	ForbiddenError      = web.NewError(1001, "已被禁言")
	UserDisabled        = web.NewError(1002, "账号已禁用")
	InObservationPeriod = web.NewError(1003, "账号尚在观察期")
	EmailNotVerified    = web.NewError(1004, "请先验证邮箱")
)
