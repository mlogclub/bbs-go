package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var UserScoreService = newUserScoreService()

func newUserScoreService() *userScoreService {
	return &userScoreService {}
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

