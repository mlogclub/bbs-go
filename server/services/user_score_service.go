package services

import (
	"bbs-go/model/constants"
	"errors"
	"strconv"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/repositories"
)

var UserScoreService = newUserScoreService()

func newUserScoreService() *userScoreService {
	return &userScoreService{}
}

type userScoreService struct {
}

func (s *userScoreService) Get(id int64) *model.UserScore {
	return repositories.UserScoreRepository.Get(simple.DB(), id)
}

func (s *userScoreService) Take(where ...interface{}) *model.UserScore {
	return repositories.UserScoreRepository.Take(simple.DB(), where...)
}

func (s *userScoreService) Find(cnd *simple.SqlCnd) []model.UserScore {
	return repositories.UserScoreRepository.Find(simple.DB(), cnd)
}

func (s *userScoreService) FindOne(cnd *simple.SqlCnd) *model.UserScore {
	return repositories.UserScoreRepository.FindOne(simple.DB(), cnd)
}

func (s *userScoreService) FindPageByParams(params *simple.QueryParams) (list []model.UserScore, paging *simple.Paging) {
	return repositories.UserScoreRepository.FindPageByParams(simple.DB(), params)
}

func (s *userScoreService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserScore, paging *simple.Paging) {
	return repositories.UserScoreRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *userScoreService) Create(t *model.UserScore) error {
	return repositories.UserScoreRepository.Create(simple.DB(), t)
}

func (s *userScoreService) Update(t *model.UserScore) error {
	return repositories.UserScoreRepository.Update(simple.DB(), t)
}

func (s *userScoreService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserScoreRepository.Updates(simple.DB(), id, columns)
}

func (s *userScoreService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserScoreRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *userScoreService) Delete(id int64) {
	repositories.UserScoreRepository.Delete(simple.DB(), id)
}

func (s *userScoreService) GetByUserId(userId int64) *model.UserScore {
	return s.FindOne(simple.NewSqlCnd().Eq("user_id", userId))
}

func (s *userScoreService) CreateOrUpdate(t *model.UserScore) error {
	if t.Id > 0 {
		return s.Update(t)
	} else {
		return s.Create(t)
	}
}

// IncrementCreateTopicScore 发帖获积分
func (s *userScoreService) IncrementPostTopicScore(topic *model.Topic) {
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostTopicScore <= 0 {
		logrus.Info("请配置发帖积分")
		return
	}
	err := s.addScore(topic.UserId, config.ScoreConfig.PostTopicScore, constants.EntityTopic,
		strconv.FormatInt(topic.Id, 10), "发表话题")
	if err != nil {
		logrus.Error(err)
	}
}

// IncrementPostCommentScore 跟帖获积分
func (s *userScoreService) IncrementPostCommentScore(comment *model.Comment) {
	// 非话题跟帖，跳过
	if comment.EntityType != constants.EntityTopic {
		return
	}
	config := SysConfigService.GetConfig()
	if config.ScoreConfig.PostCommentScore <= 0 {
		logrus.Info("请配置跟帖积分")
		return
	}
	err := s.addScore(comment.UserId, config.ScoreConfig.PostCommentScore, constants.EntityComment,
		strconv.FormatInt(comment.Id, 10), "发表跟帖")
	if err != nil {
		logrus.Error(err)
	}
}

// Increment 增加分数
func (s *userScoreService) Increment(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// Decrement 减少分数
func (s *userScoreService) Decrement(userId int64, score int, sourceType, sourceId, description string) error {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore 加分数，也可以加负数
func (s *userScoreService) addScore(userId int64, score int, sourceType, sourceId, description string) error {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	userScore := s.GetByUserId(userId)
	if userScore == nil {
		userScore = &model.UserScore{
			UserId:     userId,
			CreateTime: simple.NowTimestamp(),
		}
	}
	userScore.Score = userScore.Score + score
	userScore.UpdateTime = simple.NowTimestamp()
	if err := s.CreateOrUpdate(userScore); err != nil {
		return err
	}

	scoreType := constants.ScoreTypeIncr
	if score < 0 {
		scoreType = constants.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&model.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  simple.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.InvalidateScore(userId)
	}
	return err
}
