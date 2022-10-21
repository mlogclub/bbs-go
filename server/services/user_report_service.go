package services

import (
	"bbs-go/model"
	"bbs-go/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserReportService = newUserReportService()

func newUserReportService() *userReportService {
	return &userReportService {}
}

type userReportService struct {
}

func (s *userReportService) Get(id int64) *model.UserReport {
	return repositories.UserReportRepository.Get(sqls.DB(), id)
}

func (s *userReportService) Take(where ...interface{}) *model.UserReport {
	return repositories.UserReportRepository.Take(sqls.DB(), where...)
}

func (s *userReportService) Find(cnd *sqls.Cnd) []model.UserReport {
	return repositories.UserReportRepository.Find(sqls.DB(), cnd)
}

func (s *userReportService) FindOne(cnd *sqls.Cnd) *model.UserReport {
	return repositories.UserReportRepository.FindOne(sqls.DB(), cnd)
}

func (s *userReportService) FindPageByParams(params *params.QueryParams) (list []model.UserReport, paging *sqls.Paging) {
	return repositories.UserReportRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userReportService) FindPageByCnd(cnd *sqls.Cnd) (list []model.UserReport, paging *sqls.Paging) {
	return repositories.UserReportRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userReportService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserReportRepository.Count(sqls.DB(), cnd)
}

func (s *userReportService) Create(t *model.UserReport) error {
	return repositories.UserReportRepository.Create(sqls.DB(), t)
}

func (s *userReportService) Update(t *model.UserReport) error {
	return repositories.UserReportRepository.Update(sqls.DB(), t)
}

func (s *userReportService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserReportRepository.Updates(sqls.DB(), id, columns)
}

func (s *userReportService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserReportRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userReportService) Delete(id int64) {
	repositories.UserReportRepository.Delete(sqls.DB(), id)
}

