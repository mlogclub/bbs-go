package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var UserRoleService = newUserRoleService()

func newUserRoleService() *userRoleService {
	return &userRoleService{}
}

type userRoleService struct {
}

func (s *userRoleService) Get(id int64) *models.UserRole {
	return repositories.UserRoleRepository.Get(sqls.DB(), id)
}

func (s *userRoleService) Take(where ...interface{}) *models.UserRole {
	return repositories.UserRoleRepository.Take(sqls.DB(), where...)
}

func (s *userRoleService) Find(cnd *sqls.Cnd) []models.UserRole {
	return repositories.UserRoleRepository.Find(sqls.DB(), cnd)
}

func (s *userRoleService) FindOne(cnd *sqls.Cnd) *models.UserRole {
	return repositories.UserRoleRepository.FindOne(sqls.DB(), cnd)
}

func (s *userRoleService) FindPageByParams(params *params.QueryParams) (list []models.UserRole, paging *sqls.Paging) {
	return repositories.UserRoleRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userRoleService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserRole, paging *sqls.Paging) {
	return repositories.UserRoleRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userRoleService) Count(cnd *sqls.Cnd) int64 {
	return repositories.UserRoleRepository.Count(sqls.DB(), cnd)
}

func (s *userRoleService) Create(t *models.UserRole) error {
	return repositories.UserRoleRepository.Create(sqls.DB(), t)
}

func (s *userRoleService) Update(t *models.UserRole) error {
	return repositories.UserRoleRepository.Update(sqls.DB(), t)
}

func (s *userRoleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserRoleRepository.Updates(sqls.DB(), id, columns)
}

func (s *userRoleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserRoleRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userRoleService) Delete(id int64) {
	repositories.UserRoleRepository.Delete(sqls.DB(), id)
}
