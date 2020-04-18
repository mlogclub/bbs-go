package services

import (
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var TweetsService = newTweetsService()

func newTweetsService() *tweetsService {
	return &tweetsService{}
}

type tweetsService struct {
}

func (s *tweetsService) Get(id int64) *model.Tweets {
	return repositories.TweetsRepository.Get(simple.DB(), id)
}

func (s *tweetsService) Take(where ...interface{}) *model.Tweets {
	return repositories.TweetsRepository.Take(simple.DB(), where...)
}

func (s *tweetsService) Find(cnd *simple.SqlCnd) []model.Tweets {
	return repositories.TweetsRepository.Find(simple.DB(), cnd)
}

func (s *tweetsService) FindOne(cnd *simple.SqlCnd) *model.Tweets {
	return repositories.TweetsRepository.FindOne(simple.DB(), cnd)
}

func (s *tweetsService) FindPageByParams(params *simple.QueryParams) (list []model.Tweets, paging *simple.Paging) {
	return repositories.TweetsRepository.FindPageByParams(simple.DB(), params)
}

func (s *tweetsService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Tweets, paging *simple.Paging) {
	return repositories.TweetsRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *tweetsService) Count(cnd *simple.SqlCnd) int {
	return repositories.TweetsRepository.Count(simple.DB(), cnd)
}

func (s *tweetsService) GetTweets(cursor int64) (tweets []model.Tweets, nextCursor int64) {
	cnd := simple.NewSqlCnd().Eq("status", model.StatusOk).Desc("id").Limit(50)
	if cursor > 0 {
		cnd.Lt("id", cursor)
	}
	tweets = repositories.TweetsRepository.Find(simple.DB(), cnd)
	if len(tweets) > 0 {
		nextCursor = tweets[len(tweets)-1].Id
	} else {
		nextCursor = cursor
	}
	return
}

func (s *tweetsService) Publish(userId int64, content, imageList string) (*model.Tweets, error) {
	tweets := &model.Tweets{
		UserId:     userId,
		Content:    content,
		ImageList:  imageList,
		Status:     model.StatusOk,
		CreateTime: simple.NowTimestamp(),
	}
	if err := repositories.TweetsRepository.Create(simple.DB(), tweets); err != nil {
		return nil, err
	}
	return tweets, nil
}

func (s *tweetsService) Create(t *model.Tweets) error {
	return repositories.TweetsRepository.Create(simple.DB(), t)
}

func (s *tweetsService) Update(t *model.Tweets) error {
	return repositories.TweetsRepository.Update(simple.DB(), t)
}

func (s *tweetsService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TweetsRepository.Updates(simple.DB(), id, columns)
}

func (s *tweetsService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TweetsRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *tweetsService) Delete(id int64) {
	repositories.TweetsRepository.Delete(simple.DB(), id)
}
