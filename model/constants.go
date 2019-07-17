package model

var (
	// 系统配置
	SysConfigSiteTitle       = "site.title"
	SysConfigSiteDescription = "site.description"
	SysConfigSiteKeywords    = "site.keywords"

	// 模版全局变量
	TplCurrentUser     = "CurrentUser"
	TplSiteTitle       = "SiteTitle"
	TplSiteDescription = "SiteDescription"
	TplSiteKeywords    = "SiteKeywords"

	// 错误码
	ErrorCodeUserNameExists = 10 // 用户名已存在
	ErrorCodeEmailExists    = 11 // 邮箱已经存在
)
