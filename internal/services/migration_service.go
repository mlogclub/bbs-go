package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var MigrationService = newMigrationService()

func newMigrationService() *migrationService {
	return &migrationService{}
}

type migrationService struct {
}

func (s *migrationService) Get(id int64) *models.Migration {
	return repositories.MigrationRepository.Get(sqls.DB(), id)
}

func (s *migrationService) Take(where ...interface{}) *models.Migration {
	return repositories.MigrationRepository.Take(sqls.DB(), where...)
}

func (s *migrationService) Find(cnd *sqls.Cnd) []models.Migration {
	return repositories.MigrationRepository.Find(sqls.DB(), cnd)
}

func (s *migrationService) FindOne(cnd *sqls.Cnd) *models.Migration {
	return repositories.MigrationRepository.FindOne(sqls.DB(), cnd)
}

func (s *migrationService) FindPageByParams(params *params.QueryParams) (list []models.Migration, paging *sqls.Paging) {
	return repositories.MigrationRepository.FindPageByParams(sqls.DB(), params)
}

func (s *migrationService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Migration, paging *sqls.Paging) {
	return repositories.MigrationRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *migrationService) Count(cnd *sqls.Cnd) int64 {
	return repositories.MigrationRepository.Count(sqls.DB(), cnd)
}

func (s *migrationService) Create(t *models.Migration) error {
	return repositories.MigrationRepository.Create(sqls.DB(), t)
}

func (s *migrationService) Update(t *models.Migration) error {
	return repositories.MigrationRepository.Update(sqls.DB(), t)
}

func (s *migrationService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MigrationRepository.Updates(sqls.DB(), id, columns)
}

func (s *migrationService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MigrationRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *migrationService) Delete(id int64) {
	repositories.MigrationRepository.Delete(sqls.DB(), id)
}

func (s *migrationService) GetBy(version string) *models.Migration {
	return s.FindOne(sqls.NewCnd().Where("version = ?", version))
}
