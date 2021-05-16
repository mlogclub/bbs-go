package spam

import (
	"bbs-go/model"
	"bbs-go/pkg/common"
	"bbs-go/services"
)

type EmailVerifyStrategy struct{}

func (EmailVerifyStrategy) Name() string {
	return "EmailVerifyStrategy"
}

func (EmailVerifyStrategy) CheckTopic(user *model.User, form model.CreateTopicForm) error {
	if services.SysConfigService.IsCreateTopicEmailVerified() && !user.EmailVerified {
		return common.EmailNotVerified
	}
	return nil
}

func (EmailVerifyStrategy) CheckArticle(user *model.User, form model.CreateArticleForm) error {
	if services.SysConfigService.IsCreateArticleEmailVerified() && !user.EmailVerified {
		return common.EmailNotVerified
	}
	return nil
}

func (EmailVerifyStrategy) CheckComment(user *model.User, form model.CreateCommentForm) error {
	if services.SysConfigService.IsCreateCommentEmailVerified() && !user.EmailVerified {
		return common.EmailNotVerified
	}
	return nil
}
