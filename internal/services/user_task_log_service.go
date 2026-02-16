package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserTaskLogService = newUserTaskLogService()

func newUserTaskLogService() *userTaskLogService {
	return &userTaskLogService{}
}

type userTaskLogService struct {
}

func (s *userTaskLogService) Get(id int64) *models.UserTaskLog {
	return repositories.UserTaskLogRepository.Get(sqls.DB(), id)
}

func (s *userTaskLogService) Take(where ...interface{}) *models.UserTaskLog {
	return repositories.UserTaskLogRepository.Take(sqls.DB(), where...)
}

func (s *userTaskLogService) Find(cnd *sqls.Cnd) []models.UserTaskLog {
	return repositories.UserTaskLogRepository.Find(sqls.DB(), cnd)
}

func (s *userTaskLogService) FindOne(cnd *sqls.Cnd) *models.UserTaskLog {
	return repositories.UserTaskLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *userTaskLogService) FindPageByParams(params *params.QueryParams) (list []models.UserTaskLog, paging *sqls.Paging) {
	return repositories.UserTaskLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userTaskLogService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserTaskLog, paging *sqls.Paging) {
	return repositories.UserTaskLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userTaskLogService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserTaskLogRepository.Count(sqls.DB(), cnd)
}

func (s *userTaskLogService) Create(t *models.UserTaskLog) error {
	return repositories.UserTaskLogRepository.Create(sqls.DB(), t)
}

func (s *userTaskLogService) Update(t *models.UserTaskLog) error {
	return repositories.UserTaskLogRepository.Update(sqls.DB(), t)
}

func (s *userTaskLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserTaskLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *userTaskLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserTaskLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userTaskLogService) Delete(id int64) {
	repositories.UserTaskLogRepository.Delete(sqls.DB(), id)
}

