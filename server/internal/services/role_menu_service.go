package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/arrs"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
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

func (s *roleMenuService) GetByRole(roleId int64) []models.RoleMenu {
	return s.Find(sqls.NewCnd().Eq("role_id", roleId))
}

func (s *roleMenuService) GetMenuIdsByRoles(roleIds []int64) (menuIds []int64) {
	list := s.Find(sqls.NewCnd().In("role_id", roleIds))
	for _, element := range list {
		menuIds = append(menuIds, element.MenuId)
	}
	return
}

func (s *roleMenuService) GetMenuIdsByRole(roleId int64) (menuIds []int64) {
	list := s.GetByRole(roleId)
	for _, element := range list {
		menuIds = append(menuIds, element.MenuId)
	}
	return
}

func (s *roleMenuService) SaveRoleMenus(roleId int64, menuIds []int64) error {
	currentMenuIds := s.GetMenuIdsByRole(roleId)
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		var (
			addIds []int64 // 本次需要新增的
			delIds []int64 // 本次需要删除的
		)
		for _, menuId := range menuIds {
			if !arrs.Contains(currentMenuIds, menuId) {
				addIds = append(addIds, menuId)
			}
		}
		for _, menuId := range currentMenuIds {
			if !arrs.Contains(menuIds, menuId) {
				delIds = append(delIds, menuId)
			}
		}

		for _, menuId := range addIds {
			if err := repositories.RoleMenuRepository.Create(tx, &models.RoleMenu{
				RoleId:     roleId,
				MenuId:     menuId,
				CreateTime: dates.NowTimestamp(),
			}); err != nil {
				return err
			}
		}
		for _, menuId := range delIds {
			if err := tx.Delete(&models.RoleMenu{}, "role_id = ? and menu_id = ?", roleId, menuId).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
