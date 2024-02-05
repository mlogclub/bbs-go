package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var RoleService = newRoleService()

func newRoleService() *roleService {
	return &roleService{}
}

type roleService struct {
}

func (s *roleService) Get(id int64) *models.Role {
	return repositories.RoleRepository.Get(sqls.DB(), id)
}

func (s *roleService) Take(where ...interface{}) *models.Role {
	return repositories.RoleRepository.Take(sqls.DB(), where...)
}

func (s *roleService) Find(cnd *sqls.Cnd) []models.Role {
	return repositories.RoleRepository.Find(sqls.DB(), cnd)
}

func (s *roleService) FindOne(cnd *sqls.Cnd) *models.Role {
	return repositories.RoleRepository.FindOne(sqls.DB(), cnd)
}

func (s *roleService) FindPageByParams(params *params.QueryParams) (list []models.Role, paging *sqls.Paging) {
	return repositories.RoleRepository.FindPageByParams(sqls.DB(), params)
}

func (s *roleService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Role, paging *sqls.Paging) {
	return repositories.RoleRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *roleService) Count(cnd *sqls.Cnd) int64 {
	return repositories.RoleRepository.Count(sqls.DB(), cnd)
}

func (s *roleService) Create(t *models.Role) error {
	return repositories.RoleRepository.Create(sqls.DB(), t)
}

func (s *roleService) Update(t *models.Role) error {
	return repositories.RoleRepository.Update(sqls.DB(), t)
}

func (s *roleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.RoleRepository.Updates(sqls.DB(), id, columns)
}

func (s *roleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.RoleRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *roleService) Delete(id int64) {
	repositories.RoleRepository.Delete(sqls.DB(), id)
}
