package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var EmailLogService = newEmailLogService()

func newEmailLogService() *emailLogService {
	return &emailLogService{}
}

type emailLogService struct {
}

func (s *emailLogService) Get(id int64) *models.EmailLog {
	return repositories.EmailLogRepository.Get(sqls.DB(), id)
}

func (s *emailLogService) Take(where ...interface{}) *models.EmailLog {
	return repositories.EmailLogRepository.Take(sqls.DB(), where...)
}

func (s *emailLogService) Find(cnd *sqls.Cnd) []models.EmailLog {
	return repositories.EmailLogRepository.Find(sqls.DB(), cnd)
}

func (s *emailLogService) FindOne(cnd *sqls.Cnd) *models.EmailLog {
	return repositories.EmailLogRepository.FindOne(sqls.DB(), cnd)
}

func (s *emailLogService) FindPageByParams(params *params.QueryParams) (list []models.EmailLog, paging *sqls.Paging) {
	return repositories.EmailLogRepository.FindPageByParams(sqls.DB(), params)
}

func (s *emailLogService) FindPageByCnd(cnd *sqls.Cnd) (list []models.EmailLog, paging *sqls.Paging) {
	return repositories.EmailLogRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *emailLogService) Count(cnd *sqls.Cnd) int64 {
	return repositories.EmailLogRepository.Count(sqls.DB(), cnd)
}

func (s *emailLogService) Create(t *models.EmailLog) error {
	return repositories.EmailLogRepository.Create(sqls.DB(), t)
}

func (s *emailLogService) Update(t *models.EmailLog) error {
	return repositories.EmailLogRepository.Update(sqls.DB(), t)
}

func (s *emailLogService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.EmailLogRepository.Updates(sqls.DB(), id, columns)
}

func (s *emailLogService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.EmailLogRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *emailLogService) Delete(id int64) {
	repositories.EmailLogRepository.Delete(sqls.DB(), id)
}
