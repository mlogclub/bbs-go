package model

// 站点导航
type SiteNav struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

// 积分配置
type ScoreConfig struct {
	PostTopicScore   int `json:"postTopicScore"`   // 发帖获得积分
	PostCommentScore int `json:"postCommentScore"` // 跟帖获得积分
}

// 配置返回结构体
type ConfigData struct {
	SiteTitle        string      `json:"siteTitle"`
	SiteDescription  string      `json:"siteDescription"`
	SiteKeywords     []string    `json:"siteKeywords"`
	SiteNavs         []SiteNav   `json:"siteNavs"`
	SiteNotification string      `json:"siteNotification"`
	RecommendTags    []string    `json:"recommendTags"`
	UrlRedirect      bool        `json:"urlRedirect"`
	ScoreConfig      ScoreConfig `json:"scoreConfig"`
	DefaultNodeId    int64       `json:"defaultNodeId"`
	ArticlePending   bool        `json:"articlePending"`
}
