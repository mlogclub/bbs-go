package model

import "bbs-go/model/constants"

type UserInfo struct {
	Id                   int64    `json:"id"`
	Username             string   `json:"username"`
	Email                string   `json:"email"`
	EmailVerified        bool     `json:"emailVerified"`
	Nickname             string   `json:"nickname"`
	Avatar               string   `json:"avatar"`
	SmallAvatar          string   `json:"smallAvatar"`
	BackgroundImage      string   `json:"backgroundImage"`
	SmallBackgroundImage string   `json:"smallBackgroundImage"`
	Type                 int      `json:"type"`
	Roles                []string `json:"roles"`
	HomePage             string   `json:"homePage"`
	Description          string   `json:"description"`
	Score                int      `json:"score"`        // 积分
	TopicCount           int      `json:"topicCount"`   // 话题数量
	CommentCount         int      `json:"commentCount"` // 跟帖数量
	PasswordSet          bool     `json:"passwordSet"`  // 密码已设置
	Forbidden            bool     `json:"forbidden"`    // 是否禁言
	Status               int      `json:"status"`
	CreateTime           int64    `json:"createTime"`
}

type TagResponse struct {
	TagId   int64  `json:"tagId"`
	TagName string `json:"tagName"`
}

type ArticleSimpleResponse struct {
	ArticleId  int64          `json:"articleId"`
	User       *UserInfo      `json:"user"`
	Tags       *[]TagResponse `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	SourceUrl  string         `json:"sourceUrl"`
	ViewCount  int64          `json:"viewCount"`
	CreateTime int64          `json:"createTime"`
	Status     int            `json:"status"`
}

type ArticleResponse struct {
	ArticleSimpleResponse
	Content string `json:"content"`
}

type NodeResponse struct {
	NodeId      int64  `json:"nodeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SearchTopicResponse struct {
	TopicId    int64          `json:"topicId"`
	User       *UserInfo      `json:"user"`
	Node       *NodeResponse  `json:"node"`
	Tags       *[]TagResponse `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	CreateTime int64          `json:"createTime"`
}

// 帖子列表返回实体
type TopicResponse struct {
	TopicId         int64               `json:"topicId"`
	Type            constants.TopicType `json:"type"`
	User            *UserInfo           `json:"user"`
	Node            *NodeResponse       `json:"node"`
	Tags            *[]TagResponse      `json:"tags"`
	Title           string              `json:"title"`
	Summary         string              `json:"summary"`
	Content         string              `json:"content"`
	ImageList       []ImageInfo         `json:"imageList"`
	LastCommentTime int64               `json:"lastCommentTime"`
	ViewCount       int64               `json:"viewCount"`
	CommentCount    int64               `json:"commentCount"`
	LikeCount       int64               `json:"likeCount"`
	Liked           bool                `json:"liked"`
	CreateTime      int64               `json:"createTime"`
	Recommend       bool                `json:"recommend"`
	RecommendTime   int64               `json:"recommendTime"`
}

// 帖子列表返回实体
type TweetResponse struct {
	TweetId      int64       `json:"tweetId"`
	User         *UserInfo   `json:"user"`
	Content      string      `json:"content"`
	ImageList    []ImageInfo `json:"imageList"`
	CommentCount int64       `json:"commentCount"`
	LikeCount    int64       `json:"likeCount"`
	Status       int         `json:"status"`
	CreateTime   int64       `json:"createTime"`
}

// 项目简单返回
type ProjectSimpleResponse struct {
	ProjectId   int64     `json:"projectId"`
	User        *UserInfo `json:"user"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Logo        string    `json:"logo"`
	Url         string    `json:"url"`
	DocUrl      string    `json:"docUrl"`
	DownloadUrl string    `json:"downloadUrl"`
	Summary     string    `json:"summary"`
	CreateTime  int64     `json:"createTime"`
}

// 项目详情
type ProjectResponse struct {
	ProjectSimpleResponse
	Content string `json:"content"`
}

// 评论返回数据
type CommentResponse struct {
	CommentId    int64            `json:"commentId"`
	User         *UserInfo        `json:"user"`
	EntityType   string           `json:"entityType"`
	EntityId     int64            `json:"entityId"`
	Content      string           `json:"content"`
	QuoteId      int64            `json:"quoteId"`
	Quote        *CommentResponse `json:"quote"`
	QuoteContent string           `json:"quoteContent"`
	Status       int              `json:"status"`
	CreateTime   int64            `json:"createTime"`
}

// 收藏返回数据
type FavoriteResponse struct {
	FavoriteId int64     `json:"favoriteId"`
	EntityType string    `json:"entityType"`
	EntityId   int64     `json:"entityId"`
	Deleted    bool      `json:"deleted"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	User       *UserInfo `json:"user"`
	Url        string    `json:"url"`
	CreateTime int64     `json:"createTime"`
}

// 消息
type MessageResponse struct {
	MessageId    int64     `json:"messageId"`
	From         *UserInfo `json:"from"`    // 消息发送人
	UserId       int64     `json:"userId"`  // 消息接收人编号
	Title        string    `json:"title"`   // 标题
	Content      string    `json:"content"` // 消息内容
	QuoteContent string    `json:"quoteContent"`
	Type         int       `json:"type"`
	DetailUrl    string    `json:"detailUrl"` // 消息详情url
	ExtraData    string    `json:"extraData"`
	Status       int       `json:"status"`
	CreateTime   int64     `json:"createTime"`
}

// 图片
type ImageInfo struct {
	Url     string `json:"url"`
	Preview string `json:"preview"`
}
