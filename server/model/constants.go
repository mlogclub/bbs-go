package model

// 系统配置
const (
	SysConfigSiteTitle        = "siteTitle"        // 站点标题
	SysConfigSiteDescription  = "siteDescription"  // 站点描述
	SysConfigSiteKeywords     = "siteKeywords"     // 站点关键字
	SysConfigSiteNavs         = "siteNavs"         // 站点导航
	SysConfigSiteNotification = "siteNotification" // 站点公告
	SysConfigRecommendTags    = "recommendTags"    // 推荐标签
	SysConfigUrlRedirect      = "urlRedirect"      // 是否开启链接跳转
	SysConfigScoreConfig      = "scoreConfig"      // 分数配置
	SysConfigDefaultNodeId    = "defaultNodeId"    // 发帖默认节点
	SysConfigArticlePending   = "articlePending"   // 是否开启文章审核
	SysConfigTopicCaptcha     = "topicCaptcha"     // 是否开启发帖验证码
)

const (
	ROLE_USER    = "用户"
	ROLE_ADMIN   = "管理员"
	ROLE_MANAGER = "副站长"
)
