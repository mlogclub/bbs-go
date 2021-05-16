package spam

import "bbs-go/model"

type Strategy interface {
	// Name 策略名称
	Name() string
	// CheckTopic 检查话题
	CheckTopic(user *model.User, form model.CreateTopicForm) error
	// CheckArticle 检查文章
	CheckArticle(user *model.User, form model.CreateArticleForm) error
	// CheckComment 检查评论
	CheckComment(user *model.User, form model.CreateCommentForm) error
}
