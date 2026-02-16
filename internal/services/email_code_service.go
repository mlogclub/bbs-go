package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var EmailCodeService = newEmailCodeService()

func newEmailCodeService() *emailCodeService {
	return &emailCodeService{}
}

type emailCodeService struct {
}

func (s *emailCodeService) Get(id int64) *models.EmailCode {
	return repositories.EmailCodeRepository.Get(sqls.DB(), id)
}

func (s *emailCodeService) Take(where ...interface{}) *models.EmailCode {
	return repositories.EmailCodeRepository.Take(sqls.DB(), where...)
}

func (s *emailCodeService) Find(cnd *sqls.Cnd) []models.EmailCode {
	return repositories.EmailCodeRepository.Find(sqls.DB(), cnd)
}

func (s *emailCodeService) FindOne(cnd *sqls.Cnd) *models.EmailCode {
	return repositories.EmailCodeRepository.FindOne(sqls.DB(), cnd)
}

func (s *emailCodeService) FindPageByParams(params *params.QueryParams) (list []models.EmailCode, paging *sqls.Paging) {
	return repositories.EmailCodeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *emailCodeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.EmailCode, paging *sqls.Paging) {
	return repositories.EmailCodeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *emailCodeService) Count(cnd *sqls.Cnd) int64 {
	return repositories.EmailCodeRepository.Count(sqls.DB(), cnd)
}

func (s *emailCodeService) Create(t *models.EmailCode) error {
	return repositories.EmailCodeRepository.Create(sqls.DB(), t)
}

func (s *emailCodeService) Update(t *models.EmailCode) error {
	return repositories.EmailCodeRepository.Update(sqls.DB(), t)
}

func (s *emailCodeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.EmailCodeRepository.Updates(sqls.DB(), id, columns)
}

func (s *emailCodeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.EmailCodeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *emailCodeService) Delete(id int64) {
	repositories.EmailCodeRepository.Delete(sqls.DB(), id)
}
