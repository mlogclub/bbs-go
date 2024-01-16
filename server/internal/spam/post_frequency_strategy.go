package spam

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
	"errors"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
)

// PostFrequencyStrategy 发表频率限制
type PostFrequencyStrategy struct{}

func (PostFrequencyStrategy) Name() string {
	return "PostFrequencyStrategy"
}

func (PostFrequencyStrategy) CheckTopic(user *models.User, topic models.CreateTopicForm) error {
	// 注册时间超过24小时
	if user.CreateTime < dates.Timestamp(time.Now().Add(-time.Hour*24)) {
		return nil
	}
	var (
		maxCountInTenMinutes int64 = 1 // 十分钟内最高发帖数量
		maxCountInOneHour    int64 = 2 // 一小时内最高发帖量
		maxCountInOneDay     int64 = 3 // 一天内最高发帖量
	)

	if repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.TopicRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("发表太快了，请休息一会儿")
	}

	return nil
}

func (s PostFrequencyStrategy) CheckArticle(user *models.User, form models.CreateArticleForm) error {
	// 注册时间超过24小时
	if user.CreateTime < dates.Timestamp(time.Now().Add(-time.Hour*24)) {
		return nil
	}
	var (
		maxCountInTenMinutes int64 = 1 // 十分钟内最高发帖数量
		maxCountInOneHour    int64 = 2 // 一小时内最高发帖量
		maxCountInOneDay     int64 = 3 // 一天内最高发帖量
	)

	if repositories.ArticleRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.ArticleRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.ArticleRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("发表太快了，请休息一会儿")
	}

	return nil
}

func (s PostFrequencyStrategy) CheckComment(user *models.User, form models.CreateCommentForm) error {
	// 注册时间超过24小时
	if user.CreateTime < dates.Timestamp(time.Now().Add(-time.Hour*24)) {
		return nil
	}

	var (
		maxCountInTenMinutes int64 = 1 // 十分钟内最高发帖数量
		maxCountInOneHour    int64 = 1 // 一小时内最高发帖量
		maxCountInOneDay     int64 = 1 // 一天内最高发帖量
	)

	if repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.CommentRepository.Count(sqls.DB(), sqls.NewCnd().Eq("user_id", user.Id).
		Gt("create_time", dates.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("发表太快了，请休息一会儿")
	}
	return nil
}
