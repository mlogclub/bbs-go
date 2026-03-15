package dto

// SysConfigAdminResponse
//
//	Admin配置返回结构体
type SysConfigAdminResponse struct {
	SiteTitle                  string                      `json:"siteTitle"`
	SiteDescription            string                      `json:"siteDescription"`
	BaseURL                    string                      `json:"baseURL"`
	SiteKeywords               []string                    `json:"siteKeywords"`
	SiteLogo                   string                      `json:"siteLogo"`
	SiteNavs                   []ActionLink                `json:"siteNavs"`
	SiteNotification           string                      `json:"siteNotification"`
	AboutPageConfig            AboutPageConfig             `json:"aboutPageConfig"`
	FooterLinks                []FooterLink                `json:"footerLinks"`
	RecommendTags              []string                    `json:"recommendTags"`
	UrlRedirect                bool                        `json:"urlRedirect"`
	DefaultNodeId              int64                       `json:"defaultNodeId"`
	ArticlePending             bool                        `json:"articlePending"`
	TopicCaptcha               bool                        `json:"topicCaptcha"`
	UserObserveSeconds         int                         `json:"userObserveSeconds"`
	TokenExpireDays            int                         `json:"tokenExpireDays"`
	CreateTopicEmailVerified   bool                        `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool                        `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool                        `json:"createCommentEmailVerified"`
	EnableHideContent          bool                        `json:"enableHideContent"`
	EnableQaBounty             bool                        `json:"enableQaBounty"`
	QaBountyMin                int                         `json:"qaBountyMin"`
	QaBountyMax                int                         `json:"qaBountyMax"`
	QaBountyRequired           bool                        `json:"qaBountyRequired"`
	Modules                    ModulesConfig               `json:"modules"`
	EmailWhitelist             []string                    `json:"emailWhitelist"`             // 邮箱白名单
	EmailNoticeIntervalSeconds int                         `json:"emailNoticeIntervalSeconds"` // 邮件通知间隔(秒)
	NotificationTypes          map[string]NoticeTypeConfig `json:"notificationTypes"`          // 各消息类型站内信+邮件开关
	LoginConfig                LoginConfig                 `json:"loginConfig"`                // 登录配置
	SmtpConfig                 SmtpConfig                  `json:"smtpConfig"`                 // SMTP配置
	UploadConfig               UploadConfig                `json:"uploadConfig"`               // 上传配置
	AttachmentConfig           AttachmentConfig            `json:"attachmentConfig"`           // 附件配置
	ScriptInjections           []ScriptInjection           `json:"scriptInjections"`           // head脚本注入
}

// SysConfigOpenResponse
//
//	Open配置返回结构体
type SysConfigOpenResponse struct {
	SiteTitle                  string            `json:"siteTitle"`
	SiteDescription            string            `json:"siteDescription"`
	BaseURL                    string            `json:"baseURL"`
	SiteKeywords               []string          `json:"siteKeywords"`
	SiteLogo                   string            `json:"siteLogo"`
	SiteNavs                   []ActionLink      `json:"siteNavs"`
	SiteNotification           string            `json:"siteNotification"`
	FooterLinks                []FooterLink      `json:"footerLinks"`
	RecommendTags              []string          `json:"recommendTags"`
	UrlRedirect                bool              `json:"urlRedirect"`
	DefaultNodeId              int64             `json:"defaultNodeId"`
	ArticlePending             bool              `json:"articlePending"`
	TopicCaptcha               bool              `json:"topicCaptcha"`
	UserObserveSeconds         int               `json:"userObserveSeconds"`
	TokenExpireDays            int               `json:"tokenExpireDays"`
	CreateTopicEmailVerified   bool              `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool              `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool              `json:"createCommentEmailVerified"`
	EnableHideContent          bool              `json:"enableHideContent"`
	EnableQaBounty             bool              `json:"enableQaBounty"`
	QaBountyMin                int               `json:"qaBountyMin"`
	QaBountyMax                int               `json:"qaBountyMax"`
	QaBountyRequired           bool              `json:"qaBountyRequired"`
	Modules                    ModulesConfig     `json:"modules"`
	EmailNoticeIntervalSeconds int               `json:"emailNoticeIntervalSeconds"` // 邮件通知间隔(秒)
	AttachmentConfig           AttachmentConfig  `json:"attachmentConfig"`           // 附件配置
	LoginConfig                OpenLoginConfig   `json:"loginConfig"`                // 登录配置
	ScriptInjections           []ScriptInjection `json:"scriptInjections"`           // head脚本注入
}

type ScriptInjection struct {
	Enabled     bool   `json:"enabled"`
	ScriptName  string `json:"scriptName"` // 仅后台展示
	Type        string `json:"type"`       // external | inline
	Src         string `json:"src"`
	Code        string `json:"code"`
	Async       bool   `json:"async"`
	Defer       bool   `json:"defer"`
	Crossorigin string `json:"crossorigin"`
}

type AboutPageConfig struct {
	Content LocalizedText `json:"content"`
}

type FooterLink struct {
	Text            LocalizedText `json:"text"`
	Url             string        `json:"url"`
	OpenInNewWindow bool          `json:"openInNewWindow"`
	Visible         bool          `json:"visible"`
}

type OpenLoginConfig struct {
	PasswordLogin EnabledConfig `json:"passwordLogin"` // 密码登录
	WeixinLogin   EnabledConfig `json:"weixinLogin"`   // 微信登录
	SmsLogin      EnabledConfig `json:"smsLogin"`      // 短信登录
	GoogleLogin   EnabledConfig `json:"googleLogin"`   // Google登录
	GithubLogin   EnabledConfig `json:"githubLogin"`   // GitHub登录
}

type EnabledConfig struct {
	Enabled bool `json:"enabled"`
}

// NoticeTypeConfig 某类消息的站内信/邮件开关
type NoticeTypeConfig struct {
	Site  bool `json:"site"`
	Email bool `json:"email"`
}

// ModulesConfig
//
//	模块配置
type ModulesConfig struct {
	Tweet   bool `json:"tweet"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
}

// LoginConfig 登录配置
type LoginConfig struct {
	// 密码登录
	PasswordLogin EnabledConfig `json:"passwordLogin"`

	// 微信登录
	WeixinLogin struct {
		Enabled   bool   `json:"enabled"`
		AppId     string `json:"appId"`
		AppSecret string `json:"appSecret"`
	} `json:"weixinLogin"`

	// 短信登录
	SmsLogin struct {
		Enabled bool `json:"enabled"`
		// 短信平台
		Platform string `json:"platform"`
		// 阿里云平台配置
		Aliyun AliyunSmsConfig `json:"aliyun"`
	} `json:"smsLogin"`

	// Google登录
	GoogleLogin struct {
		Enabled      bool   `json:"enabled"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"googleLogin"`

	// GitHub登录
	GithubLogin struct {
		Enabled      bool   `json:"enabled"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	} `json:"githubLogin"`
}

type AliyunSmsConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

// IsAllDisabled 是否禁用了所有登录方式
func (c *LoginConfig) IsAllDisabled() bool {
	return !c.PasswordLogin.Enabled && !c.WeixinLogin.Enabled && !c.SmsLogin.Enabled && !c.GoogleLogin.Enabled && !c.GithubLogin.Enabled
}

type UploadMethod string

const (
	Local      UploadMethod = "Local"
	AliyunOss  UploadMethod = "AliyunOss"
	TencentCos UploadMethod = "TencentCos"
	AwsS3      UploadMethod = "AwsS3"
)

type UploadConfig struct {
	EnableUploadMethod UploadMethod           `json:"enableUploadMethod"`
	AliyunOss          AliyunOssUploadConfig  `json:"aliyunOss"`
	TencentCos         TencentCosUploadConfig `json:"tencentCos"`
	AwsS3              AwsS3UploadConfig      `json:"awsS3"`
}

type SmtpConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	SSL      bool   `json:"ssl"`
}

type AliyunOssUploadConfig struct {
	Host            string `json:"host"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	StyleSplitter   string `json:"styleSplitter"`
	StyleAvatar     string `json:"styleAvatar"`
	StylePreview    string `json:"stylePreview"`
	StyleSmall      string `json:"styleSmall"`
	StyleDetail     string `json:"styleDetail"`
}

type TencentCosUploadConfig struct {
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type AwsS3UploadConfig struct {
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

// AttachmentConfig 帖子附件配置（单 Key 存 JSON）
type AttachmentConfig struct {
	Enabled      bool     `json:"enabled"`      // 是否开启附件上传
	AllowedTypes []string `json:"allowedTypes"` // 允许的扩展名，如 [".pdf",".doc"]，空表示使用默认
	MaxSizeMB    int      `json:"maxSizeMB"`    // 单个附件大小限制(MB)，0 表示默认 10MB
	MaxCount     int      `json:"maxCount"`     // 每篇帖子最多附件数，0 表示默认 5
}
