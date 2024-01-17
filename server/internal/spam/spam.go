package spam

import (
	"bbs-go/internal/models"
	"log/slog"
)

var strategies []Strategy

func init() {
	strategies = append(strategies, &EmailVerifyStrategy{})
	strategies = append(strategies, &CaptchaStrategy{})
	// strategies = append(strategies, &PostFrequencyStrategy{})
}

func CheckTopic(user *models.User, form models.CreateTopicForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckTopic(user, form); err != nil {
			slog.Warn("[Topic]命中策略", slog.Any("strategy", strategy.Name()), slog.Any("userId", user.Id))
			return err
		}
	}
	return nil
}

func CheckArticle(user *models.User, form models.CreateArticleForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckArticle(user, form); err != nil {
			slog.Warn("[Article]命中策略", slog.Any("strategy", strategy.Name()), slog.Any("userId", user.Id))
			return err
		}
	}
	return nil
}

func CheckComment(user *models.User, form models.CreateCommentForm) error {
	if len(strategies) == 0 {
		return nil
	}
	for _, strategy := range strategies {
		if err := strategy.CheckComment(user, form); err != nil {
			slog.Warn("[Comment]命中策略", slog.Any("strategy", strategy.Name()), slog.Any("userId", user.Id))
			return err
		}
	}
	return nil
}
