package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserExpLogService = newUserExpLogService()

func newUserExpLogService() *userExpLogService {
	return &userExpLogService{}
}

type userExpLogService struct {
}

func (s *userExpLogService) Get(id int64) *models.UserExpLog {
	return repositories.UserExpLogRepository.Get(sqls.DB(), id)
}

func (s *userExpLogService) Take(where ...interface{}) *models.UserExpLog {
	return repositories.UserExpLogRepository.Take(sqls.DB(), where...)
}

func (s *userExpLogService) Find(cnd *sqls.Cnd) []models.UserExpLog {
	return repositories.UserExpLogRepository.Find(sqls.DB(), cnd)
}

func (s *userExpLogService) FindOne(cnd *sqls.Cnd) *models.UserExpLog {
	return repositories.UserExpLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *userExpLogService) FindPageByParams(params *params.QueryParams) (list []models.UserExpLog, paging *sqls.Paging) {
	return repositories.UserExpLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userExpLogService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserExpLog, paging *sqls.Paging) {
	return repositories.UserExpLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userExpLogService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserExpLogRepository.Count(sqls.DB(), cnd)
}

func (s *userExpLogService) Create(t *models.UserExpLog) error {
	return repositories.UserExpLogRepository.Create(sqls.DB(), t)
}

func (s *userExpLogService) Update(t *models.UserExpLog) error {
	return repositories.UserExpLogRepository.Update(sqls.DB(), t)
}

func (s *userExpLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserExpLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *userExpLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserExpLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userExpLogService) Delete(id int64) {
	repositories.UserExpLogRepository.Delete(sqls.DB(), id)
}

