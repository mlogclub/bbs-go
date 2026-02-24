package dto

// SysConfigAdminResponse
//
//	Admin配置返回结构体
type SysConfigAdminResponse struct {
	SiteTitle                  string        `json:"siteTitle"`
	SiteDescription            string        `json:"siteDescription"`
	BaseURL                    string        `json:"baseURL"`
	SiteKeywords               []string      `json:"siteKeywords"`
	SiteLogo                   string        `json:"siteLogo"`
	SiteNavs                   []ActionLink  `json:"siteNavs"`
	SiteNotification           string        `json:"siteNotification"`
	RecommendTags              []string      `json:"recommendTags"`
	UrlRedirect                bool          `json:"urlRedirect"`
	DefaultNodeId              int64         `json:"defaultNodeId"`
	ArticlePending             bool          `json:"articlePending"`
	TopicCaptcha               bool          `json:"topicCaptcha"`
	UserObserveSeconds         int           `json:"userObserveSeconds"`
	TokenExpireDays            int           `json:"tokenExpireDays"`
	CreateTopicEmailVerified   bool          `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool          `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool          `json:"createCommentEmailVerified"`
	EnableHideContent          bool          `json:"enableHideContent"`
	Modules                    ModulesConfig `json:"modules"`
	EmailWhitelist             []string      `json:"emailWhitelist"`             // 邮箱白名单
	EmailNoticeIntervalSeconds int           `json:"emailNoticeIntervalSeconds"` // 邮件通知间隔(秒)
	LoginConfig                LoginConfig   `json:"loginConfig"`                // 登录配置
	UploadConfig               UploadConfig  `json:"uploadConfig"`               // 上传配置
}

// SysConfigOpenResponse
//
//	Open配置返回结构体
type SysConfigOpenResponse struct {
	SiteTitle                  string          `json:"siteTitle"`
	SiteDescription            string          `json:"siteDescription"`
	BaseURL                    string          `json:"baseURL"`
	SiteKeywords               []string        `json:"siteKeywords"`
	SiteLogo                   string          `json:"siteLogo"`
	SiteNavs                   []ActionLink    `json:"siteNavs"`
	SiteNotification           string          `json:"siteNotification"`
	RecommendTags              []string        `json:"recommendTags"`
	UrlRedirect                bool            `json:"urlRedirect"`
	DefaultNodeId              int64           `json:"defaultNodeId"`
	ArticlePending             bool            `json:"articlePending"`
	TopicCaptcha               bool            `json:"topicCaptcha"`
	UserObserveSeconds         int             `json:"userObserveSeconds"`
	TokenExpireDays            int             `json:"tokenExpireDays"`
	CreateTopicEmailVerified   bool            `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool            `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool            `json:"createCommentEmailVerified"`
	EnableHideContent          bool            `json:"enableHideContent"`
	Modules                    ModulesConfig   `json:"modules"`
	EmailNoticeIntervalSeconds int             `json:"emailNoticeIntervalSeconds"` // 邮件通知间隔(秒)
	LoginConfig                OpenLoginConfig `json:"loginConfig"`                // 登录配置
}

type OpenLoginConfig struct {
	PasswordLogin EnabledConfig `json:"passwordLogin"` // 密码登录
	WeixinLogin   EnabledConfig `json:"weixinLogin"`   // 微信登录
	SmsLogin      EnabledConfig `json:"smsLogin"`      // 短信登录
	GoogleLogin   EnabledConfig `json:"googleLogin"`   // Google登录
}

type EnabledConfig struct {
	Enabled bool `json:"enabled"`
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
}

type AliyunSmsConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

// IsAllDisabled 是否禁用了所有登录方式
func (c *LoginConfig) IsAllDisabled() bool {
	return !c.PasswordLogin.Enabled && !c.WeixinLogin.Enabled && !c.SmsLogin.Enabled && !c.GoogleLogin.Enabled
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
