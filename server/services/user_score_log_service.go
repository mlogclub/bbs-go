package services

import (
	"github.com/mlogclub/simple"

	"bbs-go/model"
	"bbs-go/repositories"
)

var UserScoreLogService = newUserScoreLogService()

func newUserScoreLogService() *userScoreLogService {
	return &userScoreLogService{}
}

type userScoreLogService struct {
}

func (s *userScoreLogService) Get(id int64) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.Get(simple.DB(), id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.Take(simple.DB(), where...)
}

func (s *userScoreLogService) Find(cnd *simple.SqlCnd) []model.UserScoreLog {
	return repositories.UserScoreLogRepository.Find(simple.DB(), cnd)
}

func (s *userScoreLogService) FindOne(cnd *simple.SqlCnd) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.FindOne(simple.DB(), cnd)
}

func (s *userScoreLogService) FindPageByParams(params *simple.QueryParams) (list []model.UserScoreLog, paging *simple.Paging) {
	return repositories.UserScoreLogRepository.FindPageByParams(simple.DB(), params)
}

func (s *userScoreLogService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.UserScoreLog, paging *simple.Paging) {
	return repositories.UserScoreLogRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Create(simple.DB(), t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Update(simple.DB(), t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserScoreLogRepository.Updates(simple.DB(), id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserScoreLogRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	repositories.UserScoreLogRepository.Delete(simple.DB(), id)
}
