package spam

import "bbs-go/model"

type Strategy interface {
	Name() string
	CheckTopic(user *model.User, form model.CreateTopicForm) error
	CheckArticle(user *model.User, form model.CreateArticleForm) error
	CheckComment(user *model.User, form model.CreateCommentForm) error
}
