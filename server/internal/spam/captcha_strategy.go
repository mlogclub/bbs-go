package spam

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	"github.com/dchest/captcha"
)

type CaptchaStrategy struct{}

func (CaptchaStrategy) Name() string {
	return "CaptchaStrategy"
}

func (CaptchaStrategy) CheckTopic(user *models.User, form models.CreateTopicForm) error {
	if services.SysConfigService.GetConfig().TopicCaptcha && !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
		return errs.CaptchaError
	}
	return nil
}

func (CaptchaStrategy) CheckArticle(user *models.User, form models.CreateArticleForm) error {
	return nil
}

func (CaptchaStrategy) CheckComment(user *models.User, form models.CreateCommentForm) error {
	return nil
}
