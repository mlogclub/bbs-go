package constants

const (
	DefaultTokenExpireDays       = 7   // 用户登录token默认有效期
	SummaryLen                   = 256 // 摘要长度
	UploadMaxM                   = 10
	UploadMaxBytes         int64 = 1024 * 1024 * 1024 * UploadMaxM
	CookieTokenKey               = "bbsgo_token"
)

const (
	EmailCodeBizTypeEmailVerify   = "emailVerify"
	EmailCodeBizTypePasswordReset = "passwordReset"
)

const (
	EmailLogBizTypeUnknown       = "unknown"
	EmailLogBizTypeEmailVerify   = "emailVerify"
	EmailLogBizTypePasswordReset = "passwordReset"
	EmailLogBizTypeMessageNotice = "messageNotice"
)

const (
	EmailLogStatusSuccess = 0
	EmailLogStatusFailed  = 1
)

// 系统配置
const (
	SysConfigSiteTitle                  = "siteTitle"                  // 站点标题
	SysConfigSiteDescription            = "siteDescription"            // 站点描述
	SysConfigBaseURL                    = "baseURL"                    // 网站URL
	SysConfigSiteKeywords               = "siteKeywords"               // 站点关键字
	SysConfigSiteLogo                   = "siteLogo"                   // 站点Logo
	SysConfigSiteNavs                   = "siteNavs"                   // 站点导航
	SysConfigSiteNotification           = "siteNotification"           // 站点公告
	SysConfigRecommendTags              = "recommendTags"              // 推荐标签
	SysConfigUrlRedirect                = "urlRedirect"                // 是否开启链接跳转
	SysConfigDefaultNodeId              = "defaultNodeId"              // 发帖默认节点
	SysConfigArticlePending             = "articlePending"             // 是否开启文章审核
	SysConfigTopicCaptcha               = "topicCaptcha"               // 是否开启发帖验证码
	SysConfigUserObserveSeconds         = "userObserveSeconds"         // 新用户观察期
	SysConfigTokenExpireDays            = "tokenExpireDays"            // 登录Token有效天数
	SysConfigEnableHideContent          = "enableHideContent"          // 启用回复可见功能
	SysConfigCreateTopicEmailVerified   = "createTopicEmailVerified"   // 发话题需要邮箱认证
	SysConfigCreateArticleEmailVerified = "createArticleEmailVerified" // 发话题需要邮箱认证
	SysConfigCreateCommentEmailVerified = "createCommentEmailVerified" // 发话题需要邮箱认证
	SysConfigModules                    = "modules"                    // 功能模块
	SysConfigEmailWhitelist             = "emailWhitelist"             // 邮箱白名单
	SysConfigEmailNoticeIntervalSeconds = "emailNoticeIntervalSeconds" // 邮件通知间隔(秒)
	SysConfigLoginConfig                = "loginConfig"                // 登录配置
	SysConfigSmtpConfig                 = "smtpConfig"                 // SMTP配置
	SysConfigUploadConfig               = "uploadConfig"               // 上传配置
	SysConfigScriptInjections           = "scriptInjections"           // head脚本注入配置
	SysConfigEnableQaBounty             = "enableQaBounty"             // 是否开启问答悬赏
	SysConfigQaBountyMin                = "qaBountyMin"                // 问答悬赏积分下限
	SysConfigQaBountyMax                = "qaBountyMax"                // 问答悬赏积分上限
	SysConfigQaBountyRequired           = "qaBountyRequired"           // 问答帖是否必填悬赏
	SysConfigNotificationTypes          = "notificationTypes"          // 通知类型配置（站内信+邮件开关）
)

// EntityType
const (
	EntityArticle = "article"
	EntityTopic   = "topic"
	EntityComment = "comment"
	EntityUser    = "user"
	EntityCheckIn = "checkIn"
	EntityTask    = "task"
)

const (
	SourceTypeQaBounty       = "qa_bounty"        // SourceTypeQaBounty 积分来源：问答悬赏（UserScoreLog.SourceType）
	SourceTypeQaBountyRefund = "qa_bounty_refund" // SourceTypeQaBountyRefund 积分来源：问答悬赏退回（帖子删除且未采纳答案时退给题主）
)

// TaskEventType 任务事件类型（TaskConfig.EventType）
const (
	TaskEventTypeUserLogin      = "user.login"
	TaskEventTypeCheckIn        = "checkin"
	TaskEventTypeTopicCreate    = "topic.create"
	TaskEventTypeQaQuestion     = "qa.question.publish"
	TaskEventTypeQaAnswerAccept = "qa.answer.accept"
	TaskEventTypeCommentCreate  = "comment.create"
	TaskEventTypeFollowCreate   = "follow.create"
	TaskEventTypeFavoriteCreate = "favorite.create"
	TaskEventTypeLikeCreate     = "like.create"
	TaskEventTypeLevel10        = "level.10"
)

// TaskGroup 任务分组
type TaskGroup string

const (
	TaskGroupNewbie      TaskGroup = "newbie"      // 新手任务
	TaskGroupDaily       TaskGroup = "daily"       // 每日任务
	TaskGroupAchievement TaskGroup = "achievement" // 成就任务
)

type TaskPeriod int

const (
	TaskPeriodLifetime TaskPeriod = 0 // 一次性/终身
	TaskPeriodDaily    TaskPeriod = 1 // 每日
	TaskPeriodWeekly   TaskPeriod = 2 // 每周
	TaskPeriodMonthly  TaskPeriod = 3 // 每月
	TaskPeriodYearly   TaskPeriod = 4 // 每年
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
	StatusReview  = 2 // 待审核
)

// 用户类型
const (
	UserTypeNormal   = 0 // 普通用户
	UserTypeEmployee = 1 // 员工用户
)

// 角色类型
const (
	RoleTypeSystem = 0 // 系统角色
	RoleTypeCustom = 1 // 自定义角色
)

// 内容类型
type ContentType string

const (
	ContentTypeHtml     ContentType = "html"
	ContentTypeMarkdown ContentType = "markdown"
	ContentTypeText     ContentType = "text"
)

// 积分操作类型
const (
	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-
)

type TopicType int

const (
	TopicTypeTopic TopicType = 0 // 帖子
	TopicTypeTweet TopicType = 1 // 动态
	TopicTypeQA    TopicType = 2 // 问答
)

type TopicNodeType string

const (
	TopicNodeTypeNormal TopicNodeType = "normal"
	TopicNodeTypeQA     TopicNodeType = "qa"
)

type QaStatus string

const (
	QaStatusUnsolved QaStatus = "unsolved"
	QaStatusSolved   QaStatus = "solved"
)

func IsTweetTopicType(topicType TopicType) bool {
	return topicType == TopicTypeTweet
}

func IsPostTopicType(topicType TopicType) bool {
	return topicType == TopicTypeTopic || topicType == TopicTypeQA
}

func (t TopicNodeType) Supports(topicType TopicType) bool {
	nodeType := t
	if nodeType == "" {
		nodeType = TopicNodeTypeNormal
	}
	switch nodeType {
	case TopicNodeTypeQA:
		return topicType == TopicTypeQA
	default:
		return topicType == TopicTypeTopic || topicType == TopicTypeTweet
	}
}

type VoteType int

const (
	VoteTypeSingle   VoteType = 1 // 单选
	VoteTypeMultiple VoteType = 2 // 多选
)

type ThirdType string

const (
	ThirdTypeWeixin ThirdType = "weixin"
	ThirdTypeGoogle ThirdType = "google"
)

const (
	FollowStatusNONE   = 0
	FollowStatusFollow = 1
	FollowStatusBoth   = 2
)

const (
	NodeIdNewest    int64 = 0
	NodeIdRecommend int64 = -1
	NodeIdFollow    int64 = -2
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

type MenuType string

const (
	MenuTypeMenu = "menu" // 菜单
	MenuTypeFunc = "func" // 功能
)

// 模块
const (
	ModuleTweet   = "tweet"
	ModuleTopic   = "topic"
	ModuleArticle = "article"
)

const (
	ForbiddenWordTypeWord  = "word"
	ForbiddenWordTypeRegex = "regex"
)

// BadgeName 勋章名称常量
type BadgeName string

const (
	BadgeNameNewcomer     BadgeName = "badge_newcomer"      // 新人报到
	BadgeNameFirstPost    BadgeName = "badge_first_post"    // 初来乍到/First Post
	BadgeNameFirstComment BadgeName = "badge_first_comment" // 初次发言/First Comment
	BadgeNameAuthor       BadgeName = "badge_author"        // 内容创作者/Author
	BadgeNameHelper       BadgeName = "badge_helper"        // 热心助人/Helper
	BadgeNameStreak7      BadgeName = "badge_streak_7"      // 坚持一周/7-day Streak
	BadgeNameStreak30     BadgeName = "badge_streak_30"     // 月度打卡王/30-day Streak
	BadgeNameVeteran      BadgeName = "badge_veteran"       // 资深玩家/Veteran
)
