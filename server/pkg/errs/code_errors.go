package errs

import . "github.com/mlogclub/simple/web"

var (
	NotLogin            = NewError(1, "请先登录")
	UserNot             = NewError(2, "用户不存在")
	BadRequest          = NewError(3, "非法请求")
	CaptchaError        = NewError(1000, "验证码错误")
	ForbiddenError      = NewError(1001, "已被禁言")
	UserDisabled        = NewError(1002, "账号已禁用")
	InObservationPeriod = NewError(1003, "账号尚在观察期")
	EmailNotVerified    = NewError(1004, "请先验证邮箱")
	ScoreScarcity       = NewError(1005, "积分不足")
	BuyTopicError       = NewError(1006, "购买失败!请联系站长!")
	EmailVerified       = NewError(1007, "邮箱已验证")
	NotTypeEmail        = NewError(1008, "不支持该类型邮箱")
	EmailTimeout        = NewError(1009, "验证码过期")
	ErrScore            = NewError(1010, "积分必须为正数")
	ErrPayScore         = NewError(1011, "充值失败")
	RefereeGenCode      = NewError(2000, "推荐人注册码生成失败")
)
