package spam

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/req"
	"bbs-go/internal/pkg/errs"
	"bbs-go/internal/services"

	captcha2 "bbs-go/internal/pkg/captcha"

	"github.com/dchest/captcha"
)

type CaptchaStrategy struct{}

func (CaptchaStrategy) Name() string {
	return "CaptchaStrategy"
}

func (CaptchaStrategy) CheckTopic(user *models.User, form req.CreateTopicForm) error {
	if services.SysConfigService.IsTopicCaptcha() {
		if form.CaptchaProtocol == 2 {
			if !captcha2.Verify(form.CaptchaId, form.CaptchaCode) {
				return errs.CaptchaError()
			}
		} else {
			if !captcha.VerifyString(form.CaptchaId, form.CaptchaCode) {
				return errs.CaptchaError()
			}
		}
	}
	return nil
}

func (CaptchaStrategy) CheckArticle(user *models.User, form req.CreateArticleForm) error {
	return nil
}

func (CaptchaStrategy) CheckComment(user *models.User, form req.CreateCommentForm) error {
	return nil
}
