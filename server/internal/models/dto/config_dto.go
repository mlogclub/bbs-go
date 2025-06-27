package dto

// 积分配置
type ScoreConfig struct {
	PostTopicScore   int `json:"postTopicScore"`   // 发帖获得积分
	PostCommentScore int `json:"postCommentScore"` // 跟帖获得积分
	CheckInScore     int `json:"checkInScore"`     // 签到积分
}

// SysConfigResponse
//
//	配置返回结构体
type SysConfigResponse struct {
	SiteTitle                  string        `json:"siteTitle"`
	SiteDescription            string        `json:"siteDescription"`
	SiteKeywords               []string      `json:"siteKeywords"`
	SiteLogo                   string        `json:"siteLogo"`
	SiteNavs                   []ActionLink  `json:"siteNavs"`
	SiteNotification           string        `json:"siteNotification"`
	RecommendTags              []string      `json:"recommendTags"`
	UrlRedirect                bool          `json:"urlRedirect"`
	ScoreConfig                ScoreConfig   `json:"scoreConfig"`
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
	EmailWhitelist             []string      `json:"emailWhitelist"` // 邮箱白名单
	UploadConfig               UploadConfig  `json:"uploadConfig"`   // 上传配置
}

// ModulesConfig
//
//	模块配置
type ModulesConfig struct {
	Tweet   bool `json:"tweet"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
}

type AliyunSmsConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

type UploadMethod string

const (
	AliyunOss  UploadMethod = "AliyunOss"
	TencentCos UploadMethod = "TencentCos"
)

type UploadConfig struct {
	EnableUploadMethod UploadMethod           `json:"enableUploadMethod"`
	AliyunOss          AliyunOssUploadConfig  `json:"aliyunOss"`
	TencentCos         TencentCosUploadConfig `json:"tencentCos"`
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
