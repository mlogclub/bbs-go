package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var TweetsService = newTweetsService()

func newTweetsService() *tweetsService {
	return &tweetsService {}
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

