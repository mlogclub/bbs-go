package services

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

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
	return repositories.UserScoreLogRepository.Get(sqls.DB(), id)
}

func (s *userScoreLogService) Take(where ...interface{}) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.Take(sqls.DB(), where...)
}

func (s *userScoreLogService) Find(cnd *sqls.Cnd) []model.UserScoreLog {
	return repositories.UserScoreLogRepository.Find(sqls.DB(), cnd)
}

func (s *userScoreLogService) FindOne(cnd *sqls.Cnd) *model.UserScoreLog {
	return repositories.UserScoreLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *userScoreLogService) FindPageByParams(params *params.QueryParams) (list []model.UserScoreLog, paging *sqls.Paging) {
	return repositories.UserScoreLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userScoreLogService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserScoreLog, paging *sqls.Paging) {
	return repositories.UserScoreLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userScoreLogService) Create(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Create(sqls.DB(), t)
}

func (s *userScoreLogService) Update(t *model.UserScoreLog) error {
	return repositories.UserScoreLogRepository.Update(sqls.DB(), t)
}

func (s *userScoreLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserScoreLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *userScoreLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserScoreLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userScoreLogService) Delete(id int64) {
	repositories.UserScoreLogRepository.Delete(sqls.DB(), id)
}
