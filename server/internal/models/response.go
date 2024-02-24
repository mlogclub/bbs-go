package models

import (
	"bbs-go/internal/models/constants"
	"time"

	"github.com/mlogclub/simple/web"
)

// UserInfo 用户简单信息
type UserInfo struct {
	Id           int64            `json:"id"`
	Type         int              `json:"type"`
	Nickname     string           `json:"nickname"`
	Avatar       string           `json:"avatar"`
	SmallAvatar  string           `json:"smallAvatar"`
	Gender       constants.Gender `json:"gender"`
	Birthday     *time.Time       `json:"birthday"`
	TopicCount   int              `json:"topicCount"`   // 话题数量
	CommentCount int              `json:"commentCount"` // 跟帖数量
	FansCount    int              `json:"fansCount"`    // 粉丝数量
	FollowCount  int              `json:"followCount"`  // 关注数量
	Score        int              `json:"score"`        // 积分
	Description  string           `json:"description"`
	CreateTime   int64            `json:"createTime"`

	Forbidden bool `json:"forbidden"` // 是否禁言
	Followed  bool `json:"followed"`  // 是否关注
}

// UserDetail 用户详细信息
type UserDetail struct {
	UserInfo
	Username             string `json:"username"`
	BackgroundImage      string `json:"backgroundImage"`
	SmallBackgroundImage string `json:"smallBackgroundImage"`
	HomePage             string `json:"homePage"`
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
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ArticleSimpleResponse struct {
	Id           int64          `json:"id"`
	User         *UserInfo      `json:"user"`
	Tags         *[]TagResponse `json:"tags"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	Cover        *ImageInfo     `json:"cover"`
	SourceUrl    string         `json:"sourceUrl"`
	ViewCount    int64          `json:"viewCount"`
	CommentCount int64          `json:"commentCount"`
	LikeCount    int64          `json:"likeCount"`
	CreateTime   int64          `json:"createTime"`
	Status       int            `json:"status"`
	Favorited    bool           `json:"favorited"`
}

type ArticleResponse struct {
	ArticleSimpleResponse
	Content string `json:"content"`
}

type NodeResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
}

type SearchTopicResponse struct {
	Id         int64          `json:"id"`
	User       *UserInfo      `json:"user"`
	Node       *NodeResponse  `json:"node"`
	Tags       *[]TagResponse `json:"tags"`
	Title      string         `json:"title"`
	Summary    string         `json:"summary"`
	CreateTime int64          `json:"createTime"`
}

// 帖子列表返回实体
type TopicResponse struct {
	Id              int64               `json:"id"`
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
	Status          int                 `json:"status"`
	Favorited       bool                `json:"favorited"`
	IpLocation      string              `json:"ipLocation"`
}

// CommentResponse 评论返回数据
type CommentResponse struct {
	Id           int64             `json:"id"`
	User         *UserInfo         `json:"user"`
	EntityType   string            `json:"entityType"`
	EntityId     int64             `json:"entityId"`
	ContentType  string            `json:"contentType"`
	Content      string            `json:"content"`
	ImageList    []ImageInfo       `json:"imageList"`
	LikeCount    int64             `json:"likeCount"`
	CommentCount int64             `json:"commentCount"`
	Liked        bool              `json:"liked"`
	QuoteId      int64             `json:"quoteId"`
	Quote        *CommentResponse  `json:"quote"`
	Replies      *web.CursorResult `json:"replies"`
	IpLocation   string            `json:"ipLocation"`
	Status       int               `json:"status"`
	CreateTime   int64             `json:"createTime"`
}

// 收藏返回数据
type FavoriteResponse struct {
	Id         int64     `json:"id"`
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
	Id           int64     `json:"id"`
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

type TreeNode struct {
	Id       int64      `json:"id"`
	Key      int64      `json:"key"`
	Title    string     `json:"title"`
	Children []TreeNode `json:"children"`
}

type MenuResponse struct {
	Id         int64  `json:"id"`
	ParentId   *int64 `json:"parentId"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	SortNo     int    `json:"sortNo"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type MenuTreeResponse struct {
	MenuResponse
	Level    int                `json:"level"`
	Children []MenuTreeResponse `json:"children"`
}
