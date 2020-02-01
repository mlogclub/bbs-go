
package services

import (
	"bbs-go/model"
	"bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var UserScoreLogService = newUserScoreLogService()

func newUserScoreLogService() *userScoreLogService {
	return &userScoreLogService {}
}

type userScoreLogService struct {
}

func (this *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.Get(simple.DB(), id)
}

func (this *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.Take(simple.DB(), where...)
}

func (this *userScoreLogService) Find(cnd *simple.SqlCnd) []model.UserScoreLog {
	return repositories.UserScoreLogRepository.Find(simple.DB(), cnd)
}

func (this *userScoreLogService) FindOne(cnd *simple.SqlCnd) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.FindOne(simple.DB(), cnd)
}

func (this *userScoreLogService) FindPageByParams(params *simple.QueryParams) (list []model.UserScoreLog, paging *simple.Paging) {
	return repositories.UserScoreLogRepository.FindPageByParams(simple.DB(), params)
}

func (this *userScoreLogService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserScoreLog, paging *simple.Paging) {
	return repositories.UserScoreLogRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *userScoreLogService) Create(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Create(simple.DB(), t)
}

func (this *userScoreLogService) Update(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Update(simple.DB(), t)
}

func (this *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserScoreLogRepository.Updates(simple.DB(), id, columns)
}

func (this *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserScoreLogRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *userScoreLogService) Delete(id int64) {
	repositories.UserScoreLogRepository.Delete(simple.DB(), id)
}

