
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

func (this *userScoreService) Get(id int64) *model.UserScore {
	return repositories.UserScoreRepository.Get(simple.DB(), id)
}

func (this *userScoreService) Take(where ...interface{}) *model.UserScore {
	return repositories.UserScoreRepository.Take(simple.DB(), where...)
}

func (this *userScoreService) Find(cnd *simple.SqlCnd) []model.UserScore {
	return repositories.UserScoreRepository.Find(simple.DB(), cnd)
}

func (this *userScoreService) FindOne(cnd *simple.SqlCnd) *model.UserScore {
	return repositories.UserScoreRepository.FindOne(simple.DB(), cnd)
}

func (this *userScoreService) FindPageByParams(params *simple.QueryParams) (list []model.UserScore, paging *simple.Paging) {
	return repositories.UserScoreRepository.FindPageByParams(simple.DB(), params)
}

func (this *userScoreService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserScore, paging *simple.Paging) {
	return repositories.UserScoreRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *userScoreService) Create(t *model.UserScore) error {
	return repositories.UserScoreRepository.Create(simple.DB(), t)
}

func (this *userScoreService) Update(t *model.UserScore) error {
	return repositories.UserScoreRepository.Update(simple.DB(), t)
}

func (this *userScoreService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserScoreRepository.Updates(simple.DB(), id, columns)
}

func (this *userScoreService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserScoreRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *userScoreService) Delete(id int64) {
	repositories.UserScoreRepository.Delete(simple.DB(), id)
}

