package model

import (
	"bbs-go/model/constants"

	"github.com/mlogclub/simple/web"
)

// UserInfo 用户简单信息
type UserInfo struct {
	Id           int64  `json:"id"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	SmallAvatar  string `json:"smallAvatar"`
	TopicCount   int    `json:"topicCount"`   // 话题数量
	CommentCount int    `json:"commentCount"` // 跟帖数量
	FansCount    int    `json:"fansCount"`    // 粉丝数量
	FollowCount  int    `json:"followCount"`  // 关注数量
	Score        int    `json:"score"`        // 积分
	Description  string `json:"description"`
	CreateTime   int64  `json:"createTime"`

	Followed bool `json:"followed"`
}

// UserDetail 用户详细信息
type UserDetail struct {
	UserInfo
	Username             string `json:"username"`
	BackgroundImage      string `json:"backgroundImage"`
	SmallBackgroundImage string `json:"smallBackgroundImage"`
	HomePage             string `json:"homePage"`
	Forbidden            bool   `json:"forbidden"` // 是否禁言
	Status               int    `json:"status"`
}

// UserProfile 用户个人信息
type UserProfile struct {
	UserDetail
	Roles         []string `json:"roles"`
	PasswordSet   bool     `json:"passwordSet"` // 密码已设置
	Email         string   `json:"email"`
	EmailVerified bool     `json:"emailVerified"`
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
	Logo        string `json:"logo"`
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
	Sticky          bool                `json:"sticky"`
	StickyTime      int64               `json:"stickyTime"`
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

// CommentResponse 评论返回数据
type CommentResponse struct {
	CommentId    int64             `json:"commentId"`
	User         *UserInfo         `json:"user"`
	EntityType   string            `json:"entityType"`
	EntityId     int64             `json:"entityId"`
	Content      string            `json:"content"`
	ImageList    []ImageInfo       `json:"imageList"`
	LikeCount    int64             `json:"likeCount"`
	CommentCount int64             `json:"commentCount"`
	Liked        bool              `json:"liked"`
	QuoteId      int64             `json:"quoteId"`
	Quote        *CommentResponse  `json:"quote"`
	Replies      *web.CursorResult `json:"replies"`
	Status       int               `json:"status"`
	CreateTime   int64             `json:"createTime"`
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
