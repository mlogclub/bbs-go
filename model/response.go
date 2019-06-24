package model

import (
	"html/template"
)

type UserResponse struct {
	Id          int64    `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Email       string   `json:"email"`
	Type        int      `json:"type"`
	Roles       []string `json:"roles"`
	Description string   `json:"description"`
	CreateTime  int64    `json:"createTime"`
}

func (this *UserResponse) HasRole(role string) bool {
	if len(this.Roles) == 0 {
		return false
	}
	for _, r := range this.Roles {
		if r == role {
			return true
		}
	}
	return false
}

type CategoryResponse struct {
	CategoryId   int64         `json:"categoryId"`
	CategoryName string        `json:"categoryName"`
	Tags         []TagResponse `json:"tags"`
}

type TagResponse struct {
	TagId   int64  `json:"tagId"`
	TagName string `json:"tagName"`
}

type ArticleResponse struct {
	ArticleId  int64             `json:"articleId"`
	User       *UserResponse     `json:"user"`
	Category   *CategoryResponse `json:"category"`
	Tags       *[]TagResponse    `json:"tags"`
	Title      string            `json:"title"`
	Summary    string            `json:"summary"`
	Content    template.HTML     `json:"content"`
	SourceUrl  string            `json:"sourceUrl"`
	CreateTime int64             `json:"createTime"`
}

type TopicResponse struct {
	TopicId         int64          `json:"topicId"`
	User            *UserResponse  `json:"user"`
	Tags            *[]TagResponse `json:"tags"`
	Title           string         `json:"title"`
	Content         template.HTML  `json:"content"`
	LastCommentTime int64          `json:"lastCommentTime"`
	ViewCount       int64          `json:"viewCount"`
	CreateTime      int64          `json:"createTime"`
}

type CommentResponse struct {
	CommentId    int64            `json:"commentId"`
	User         *UserResponse    `json:"user"`
	EntityType   string           `json:"entityType"`
	EntityId     int64            `json:"entityId"`
	Content      template.HTML    `json:"content"`
	QuoteId      int64            `json:"quoteId"`
	Quote        *CommentResponse `json:"quote"`
	QuoteContent template.HTML    `json:"quoteContent"`
	Status       int              `json:"status"`
	CreateTime   int64            `json:"createTime"`
}

type FavoriteResponse struct {
	FavoriteId int64         `json:"favoriteId"`
	EntityType string        `json:"entityType"`
	EntityId   int64         `json:"entityId"`
	Deleted    bool          `json:"deleted"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	User       *UserResponse `json:"user"`
	Url        string        `json:"url"`
	CreateTime int64         `json:"createTime"`
}

// 消息
type MessageResponse struct {
	MessageId    int64         `json:"messageId"`
	UserId       int64         `json:"userId"`
	User         *UserResponse `json:"user"`
	Content      string        `json:"content"`
	QuoteContent string        `json:"quoteContent"`
	Type         int           `json:"type"`
	ExtraData    string        `json:"extraData"`
	Status       int           `json:"status"`
	CreateTime   int64         `json:"createTime"`
}
