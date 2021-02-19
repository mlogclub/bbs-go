package constants

const (
	DefaultTokenExpireDays = 7   // 用户登录token默认有效期
	SummaryLen             = 256 // 摘要长度
)

// 系统配置
const (
	SysConfigSiteTitle          = "siteTitle"          // 站点标题
	SysConfigSiteDescription    = "siteDescription"    // 站点描述
	SysConfigSiteKeywords       = "siteKeywords"       // 站点关键字
	SysConfigSiteNavs           = "siteNavs"           // 站点导航
	SysConfigSiteNotification   = "siteNotification"   // 站点公告
	SysConfigRecommendTags      = "recommendTags"      // 推荐标签
	SysConfigUrlRedirect        = "urlRedirect"        // 是否开启链接跳转
	SysConfigScoreConfig        = "scoreConfig"        // 分数配置
	SysConfigDefaultNodeId      = "defaultNodeId"      // 发帖默认节点
	SysConfigArticlePending     = "articlePending"     // 是否开启文章审核
	SysConfigTopicCaptcha       = "topicCaptcha"       // 是否开启发帖验证码
	SysConfigUserObserveSeconds = "userObserveSeconds" // 新用户观察期
	SysConfigTokenExpireDays    = "tokenExpireDays"    // 登录Token有效天数
	SysConfigLoginMethod        = "loginMethod"        // 登录方式
)

// EntityType
const (
	EntityArticle = "article"
	EntityTopic   = "topic"
	EntityComment = "comment"
	EntityUser    = "user"
	EntityCheckIn = "checkIn"
)

// 用户角色
const (
	RoleOwner = "owner" // 站长
	RoleAdmin = "admin" // 管理员
	RoleUser  = "user"  // 用户
)

// 操作类型
const (
	OpTypeCreate          = "create"
	OpTypeDelete          = "delete"
	OpTypeUpdate          = "update"
	OpTypeForbidden       = "forbidden"
	OpTypeRemoveForbidden = "removeForbidden"
)

// 状态
const (
	StatusOk      = 0 // 正常
	StatusDeleted = 1 // 删除
	StatusPending = 2 // 待审核
)

// 用户类型
const (
	UserTypeNormal = 0 // 普通用户
	UserTypeGzh    = 1 // 公众号用户
)

// 内容类型
const (
	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"
	ContentTypeText     = "text"
)

// 消息状态
const (
	MsgStatusUnread   = 0 // 消息未读
	MsgStatusHaveRead = 1 // 消息已读
)

// 消息类型
const (
	MsgTypeTopicComment   = 0 // 收到话题评论
	MsgTypeCommentReply   = 1 // 收到他人回复
	MsgTypeTopicLike      = 2 // 收到点赞
	MsgTypeTopicFavorite  = 3 // 话题被收藏
	MsgTypeTopicRecommend = 4 // 话题被设为推荐
	MsgTypeTopicDelete    = 5 // 话题被删除
)

// 第三方账号类型
const (
	ThirdAccountTypeGithub = "github"
	ThirdAccountTypeQQ     = "qq"
)

// 积分操作类型
const (
	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-
)

type TopicType int

const (
	TopicTypeTopic TopicType = 0
	TopicTypeTweet TopicType = 1
)

type LoginMethod string

const (
	LoginMethodQQ       LoginMethod = "qq"
	LoginMethodGithub   LoginMethod = "github"
	LoginMethodPassword LoginMethod = "password"
)
