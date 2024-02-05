package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var RoleMenuService = newRoleMenuService()

func newRoleMenuService() *roleMenuService {
	return &roleMenuService{}
}

type roleMenuService struct {
}

func (s *roleMenuService) Get(id int64) *models.RoleMenu {
	return repositories.RoleMenuRepository.Get(sqls.DB(), id)
}

func (s *roleMenuService) Take(where ...interface{}) *models.RoleMenu {
	return repositories.RoleMenuRepository.Take(sqls.DB(), where...)
}

func (s *roleMenuService) Find(cnd *sqls.Cnd) []models.RoleMenu {
	return repositories.RoleMenuRepository.Find(sqls.DB(), cnd)
}

func (s *roleMenuService) FindOne(cnd *sqls.Cnd) *models.RoleMenu {
	return repositories.RoleMenuRepository.FindOne(sqls.DB(), cnd)
}

func (s *roleMenuService) FindPageByParams(params *params.QueryParams) (list []models.RoleMenu, paging *sqls.Paging) {
	return repositories.RoleMenuRepository.FindPageByParams(sqls.DB(), params)
}

func (s *roleMenuService) FindPageByCnd(cnd *sqls.Cnd) (list []models.RoleMenu, paging *sqls.Paging) {
	return repositories.RoleMenuRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *roleMenuService) Count(cnd *sqls.Cnd) int64 {
	return repositories.RoleMenuRepository.Count(sqls.DB(), cnd)
}

func (s *roleMenuService) Create(t *models.RoleMenu) error {
	return repositories.RoleMenuRepository.Create(sqls.DB(), t)
}

func (s *roleMenuService) Update(t *models.RoleMenu) error {
	return repositories.RoleMenuRepository.Update(sqls.DB(), t)
}

func (s *roleMenuService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.RoleMenuRepository.Updates(sqls.DB(), id, columns)
}

func (s *roleMenuService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.RoleMenuRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *roleMenuService) Delete(id int64) {
	repositories.RoleMenuRepository.Delete(sqls.DB(), id)
}
