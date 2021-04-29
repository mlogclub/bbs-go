package spam

import (
	"bbs-go/model"
	"github.com/sirupsen/logrus"
)

var strategies []Strategy

func init() {
	strategies = append(strategies, &EmailVerifyStrategy{})
	strategies = append(strategies, &CaptchaStrategy{})
	strategies = append(strategies, &PostFrequencyStrategy{})
}

func CheckTopic(user *model.User, form model.CreateTopicForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckTopic(user, form); err != nil {
			logrus.Warnf("[Topic]命中策略：%s, userId：%d", strategy.Name(), user.Id)
			return err
		}
	}
	return nil
}

func CheckArticle(user *model.User, form model.CreateArticleForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckArticle(user, form); err != nil {
			logrus.Warnf("[Article]命中策略：%s, userId：%d", strategy.Name(), user.Id)
			return err
		}
	}
	return nil
}

func CheckComment(user *model.User, form model.CreateCommentForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckComment(user, form); err != nil {
			logrus.Warnf("[Comment]命中策略：%s, userId：%d", strategy.Name(), user.Id)
			return err
		}
	}
	return nil
}
