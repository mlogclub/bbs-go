package constants

const (
	DefaultTokenExpireDays       = 7   // 用户登录token默认有效期
	SummaryLen                   = 256 // 摘要长度
	UploadMaxM                   = 10
	UploadMaxBytes         int64 = 1024 * 1024 * 1024 * UploadMaxM
	CookieTokenKey               = "bbsgo_token"
)

// 昵称长度限制
const (
	NicknameMinLengthZhCN = 2
	NicknameMaxLengthZhCN = 12
	NicknameMinLengthEnUS = 2
	NicknameMaxLengthEnUS = 20
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
	SysConfigAboutPageConfig            = "aboutPageConfig"            // 关于页配置
	SysConfigFooterLinks                = "footerLinks"                // 底部链接
	SysConfigRecommendTags              = "recommendTags"              // 推荐标签
	SysConfigUrlRedirect                = "urlRedirect"                // 是否开启链接跳转
	SysConfigDefaultCategoryId          = "defaultCategoryId"          // 发帖默认节点
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
	SysConfigAttachmentConfig           = "attachmentConfig"           // 附件配置（帖子附件）
	SysConfigScriptInjections           = "scriptInjections"           // head脚本注入配置
	SysConfigEnableQaBounty             = "enableQaBounty"             // 是否开启问答悬赏
	SysConfigQaBountyMin                = "qaBountyMin"                // 问答悬赏积分下限
	SysConfigQaBountyMax                = "qaBountyMax"                // 问答悬赏积分上限
	SysConfigQaBountyRequired           = "qaBountyRequired"           // 问答帖是否必填悬赏
	SysConfigNotificationTypes          = "notificationTypes"          // 通知类型配置（站内信+邮件开关）
)

// EntityType
const (
	EntityArticle    = "article"
	EntityTopic      = "topic"
	EntityComment    = "comment"
	EntityUser       = "user"
	EntityCheckIn    = "checkIn"
	EntityTask       = "task"
	EntityAttachment = "attachment"
)

const (
	SourceTypeQaBounty           = "qa_bounty"           // SourceTypeQaBounty 积分来源：问答悬赏（UserScoreLog.SourceType）
	SourceTypeQaBountyRefund     = "qa_bounty_refund"    // SourceTypeQaBountyRefund 积分来源：问答悬赏退回（帖子删除且未采纳答案时退给题主）
	SourceTypeAttachmentDownload = "attachment_download" // 附件下载扣积分
	SourceTypeAttachmentIncome   = "attachment_income"   // 附件下载帖主收入（可选）
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

type CategoryType string

const (
	CategoryTypeNormal CategoryType = "normal"
	CategoryTypeQA     CategoryType = "qa"
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

func (t CategoryType) Supports(topicType TopicType) bool {
	categoryType := t
	if categoryType == "" {
		categoryType = CategoryTypeNormal
	}
	switch categoryType {
	case CategoryTypeQA:
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
	ThirdTypeGithub ThirdType = "github"
)

const (
	FollowStatusNONE   = 0
	FollowStatusFollow = 1
	FollowStatusBoth   = 2
)

const (
	CategoryIdNewest    int64 = 0
	CategoryIdRecommend int64 = -1
	CategoryIdFollow    int64 = -2
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
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

// 热度点系统常量
const (
	HeatPointsGenesisAirdrop   = 50   // 创世空投每人 50 点
	HeatPointsDailyCheckinRate = 0.05 // 每日签到发放比例 0.05%
	HeatPointsStakeQuotaDaily  = 3    // 每人每日质押次数
	HeatPointsStakeMinAmount   = 1    // 每次质押最低额度
	HeatPointsDecayRate        = 0.02 // 每日衰减率 2%
	HeatPointsForfeitDays      = 60   // 60 天完全归零
	HeatPointsCooldownHours    = 24   // 赎回冷却期（小时）
	HeatPointsDefaultCirculation = 10000 // 活跃流通默认值（快照缺失时使用）
	
	// 阶段阈值（活跃流通百分比）
	HeatPhase1Threshold = 0.005 // 0.5% 冷帖期上限
	HeatPhase2Threshold = 0.03  // 3% 热门期阈值
	
	// 火焰等级阈值
	HeatFlameLevel2Threshold = 0.001 // 0.1%
	HeatFlameLevel3Threshold = 0.005 // 0.5%
	HeatFlameLevel4Threshold = 0.01  // 1%
	HeatFlameLevel5Threshold = 0.03  // 3%
	
	// 单人单帖上限
	HeatSingleTopicUserLimitRatio = 0.003 // 0.3%
)

// 质押状态
const (
	StakeStatusActive        = 0 // 质押中
	StakeStatusRedeemed      = 1 // 已赎回
	StakeStatusFailed        = 2 // 已结算失败（本金归零）
	StakeStatusAdminRedeemed = 3 // 管理员强制赎回（区分正常赎回，用于审计）
)

// 热度点流水类型
const (
	HeatLogTypeDailyCheckIn      = "DailyCheckIn"      // 每日签到
	HeatLogTypeStakeOut          = "StakeOut"          // 质押支出
	HeatLogTypeSettleProfit      = "SettleProfit"      // 结算收益
	HeatLogTypeSettleProfitPartial = "SettleProfitPartial" // 结算收益（削减后）
	HeatLogTypeSettleLoss        = "SettleLoss"        // 结算亏损
	HeatLogTypeRedeem            = "Redeem"            // 赎回
	HeatLogTypeDecay             = "Decay"             // 衰减
	HeatLogTypeDecayTruncated    = "DecayTruncated"    // 衰减截断
	HeatLogTypeFullForfeit       = "FullForfeit"       // 60 天完全归零
	HeatLogTypeRankReward        = "RankReward"        // 排名奖励
	HeatLogTypeAdminForceRedeem  = "AdminForceRedeem"  // 管理员强制赎回
	HeatLogTypeAccountClosure    = "AccountClosure"    // 销户
)

// 公共奖池来源
const (
	HeatPoolSourceDecayInflow    = "DecayInflow"    // 衰减流入
	HeatPoolSourceStakeLoss      = "StakeLoss"      // 质押亏损
	HeatPoolSourceFullForfeit    = "FullForfeit"    // 60 天归零
	HeatPoolSourceGenesisMint    = "GenesisMint"    // 创世铸币
	HeatPoolSourceSubsidyMint    = "SubsidyMint"    // 补贴期铸币
	HeatPoolSourceCheckInPayout  = "CheckInPayout"  // 签到发放
	HeatPoolSourceSettlePayout   = "SettlePayout"   // 结算发放
	HeatPoolSourceRankReward     = "RankReward"     // 排名奖励
	HeatPoolSourceAccountClosure = "AccountClosure" // 销户回收
)

// 铸币类型
const (
	HeatMintTypeGenesisAirdrop = "GenesisAirdrop" // 创世空投
	HeatMintTypeSubsidyMint    = "SubsidyMint"    // 补贴期铸币
)

// 结算任务类型
const (
	SettlementTaskTypeSettlement = "Settlement" // 每日结算
	SettlementTaskTypeSnapshot   = "Snapshot"   // 快照
	SettlementTaskTypeFlameRefresh = "FlameRefresh" // 火焰等级刷新
)

// 结算任务状态
const (
	SettlementTaskStatusRunning = 0 // 执行中
	SettlementTaskStatusCompleted = 1 // 已完成
	SettlementTaskStatusFailed  = 2 // 失败
)
