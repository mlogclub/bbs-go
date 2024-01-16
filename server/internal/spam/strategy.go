package spam

import "bbs-go/internal/models"

type Strategy interface {
	// Name 策略名称
	Name() string
	// CheckTopic 检查话题
	CheckTopic(user *models.User, form models.CreateTopicForm) error
	// CheckArticle 检查文章
	CheckArticle(user *models.User, form models.CreateArticleForm) error
	// CheckComment 检查评论
	CheckComment(user *models.User, form models.CreateCommentForm) error
}
