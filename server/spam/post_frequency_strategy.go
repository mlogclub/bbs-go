package spam

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"errors"
	"github.com/mlogclub/simple"
	"github.com/mlogclub/simple/date"
	"time"
)

// PostFrequencyStrategy 发表频率限制
type PostFrequencyStrategy struct{}

func (PostFrequencyStrategy) Name() string {
	return "PostFrequencyStrategy"
}

func (PostFrequencyStrategy) IsOpen() bool {
	// TODO 这个要做成配置
	return true
}

func (PostFrequencyStrategy) CheckTopic(user *model.User, topic model.CreateTopicForm) error {
	var (
		maxCountInTenMinutes int64 = 2  // 十分钟内最高发帖数量
		maxCountInOneHour    int64 = 6  // 一小时内最高发帖量
		maxCountInOneDay     int64 = 32 // 一天内最高发帖量
	)

	if user.CreateTime > date.Timestamp(time.Now().Add(-time.Hour*24)) { // 注册时间在24小时内
		maxCountInTenMinutes = 1
		maxCountInOneHour = 1
		maxCountInOneDay = 1
	}
	if repositories.TopicRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.TopicRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.TopicRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("发表太快了，请休息一会儿")
	}

	return nil
}

func (s PostFrequencyStrategy) CheckArticle(user *model.User, form model.CreateArticleForm) error {
	var (
		maxCountInTenMinutes int64 = 2  // 十分钟内最高发帖数量
		maxCountInOneHour    int64 = 6  // 一小时内最高发帖量
		maxCountInOneDay     int64 = 32 // 一天内最高发帖量
	)

	if user.CreateTime > date.Timestamp(time.Now().Add(-time.Hour*24)) { // 注册时间在24小时内
		maxCountInTenMinutes = 1
		maxCountInOneHour = 1
		maxCountInOneDay = 1
	}
	if repositories.ArticleRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.ArticleRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
		return errors.New("发表太快了，请休息一会儿")
	}

	if repositories.ArticleRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
		Gt("create_time", date.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
		return errors.New("发表太快了，请休息一会儿")
	}

	return nil
}

func (s PostFrequencyStrategy) CheckComment(user *model.User, form model.CreateCommentForm) error {
	if user.CreateTime > date.Timestamp(time.Now().Add(-time.Hour*24)) { // 注册时间在24小时内
		var (
			maxCountInTenMinutes int64 = 1 // 十分钟内最高发帖数量
			maxCountInOneHour    int64 = 1 // 一小时内最高发帖量
			maxCountInOneDay     int64 = 1 // 一天内最高发帖量
		)

		if repositories.CommentRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
			Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour*24)))) >= maxCountInOneDay {
			return errors.New("发表太快了，请休息一会儿")
		}

		if repositories.CommentRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
			Gt("create_time", date.Timestamp(time.Now().Add(-time.Hour)))) >= maxCountInOneHour {
			return errors.New("发表太快了，请休息一会儿")
		}

		if repositories.CommentRepository.Count(simple.DB(), simple.NewSqlCnd().Eq("user_id", user.Id).
			Gt("create_time", date.Timestamp(time.Now().Add(-time.Minute*10)))) >= maxCountInTenMinutes {
			return errors.New("发表太快了，请休息一会儿")
		}
	}

	return nil
}
