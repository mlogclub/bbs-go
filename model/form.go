package model

// 发表评论
type CreateCommentForm struct {
	EntityType string `form:"entityType"`
	EntityId   int64  `form:"entityId"`
	Content    string `form:"content"`
	QuoteId    int64  `form:"quoteId"`
}
