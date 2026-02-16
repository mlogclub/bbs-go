package spam

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/req"
)

type Strategy interface {
	// Name 策略名称
	Name() string
	// CheckTopic 检查话题
	CheckTopic(user *models.User, form req.CreateTopicForm) error
	// CheckArticle 检查文章
	CheckArticle(user *models.User, form req.CreateArticleForm) error
	// CheckComment 检查评论
	CheckComment(user *models.User, form req.CreateCommentForm) error
}
