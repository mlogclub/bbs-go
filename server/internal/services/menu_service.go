package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var MenuService = newMenuService()

func newMenuService() *menuService {
	return &menuService{}
}

type menuService struct {
}

func (s *menuService) Get(id int64) *models.Menu {
	return repositories.MenuRepository.Get(sqls.DB(), id)
}

func (s *menuService) Take(where ...interface{}) *models.Menu {
	return repositories.MenuRepository.Take(sqls.DB(), where...)
}

func (s *menuService) Find(cnd *sqls.Cnd) []models.Menu {
	return repositories.MenuRepository.Find(sqls.DB(), cnd)
}

func (s *menuService) FindOne(cnd *sqls.Cnd) *models.Menu {
	return repositories.MenuRepository.FindOne(sqls.DB(), cnd)
}

func (s *menuService) FindPageByParams(params *params.QueryParams) (list []models.Menu, paging *sqls.Paging) {
	return repositories.MenuRepository.FindPageByParams(sqls.DB(), params)
}

func (s *menuService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Menu, paging *sqls.Paging) {
	return repositories.MenuRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *menuService) Count(cnd *sqls.Cnd) int64 {
	return repositories.MenuRepository.Count(sqls.DB(), cnd)
}

func (s *menuService) Create(t *models.Menu) error {
	return repositories.MenuRepository.Create(sqls.DB(), t)
}

func (s *menuService) Update(t *models.Menu) error {
	return repositories.MenuRepository.Update(sqls.DB(), t)
}

func (s *menuService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MenuRepository.Updates(sqls.DB(), id, columns)
}

func (s *menuService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MenuRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *menuService) Delete(id int64) {
	repositories.MenuRepository.Delete(sqls.DB(), id)
}

func (s *menuService) GetNextSortNo(parentId int64) int {
	if max := s.FindOne(sqls.NewCnd().Eq("parent_id", parentId).Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *menuService) GetUserMenus(user *models.User) []models.Menu {
	// roleIds := UserRoleService.GetUserRoleIds(user.Id)
	// if len(roleIds) == 0 {
	// 	return nil
	// }
	// menuIds := RoleMenuService.GetMenuIdsByRoles(roleIds)
	// if len(menuIds) == 0 {
	// 	return nil
	// }
	// return repositories.MenuRepository.Find(sqls.DB(), sqls.NewCnd().In("id", menuIds).Asc("sort_no").Desc("id"))

	return repositories.MenuRepository.Find(sqls.DB(), sqls.NewCnd().Asc("sort_no").Desc("id"))
}

func (s *menuService) UpdateSort(ids []int64) error {
	return sqls.DB().Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := repositories.MenuRepository.UpdateColumn(tx, id, "sort_no", i); err != nil {
				return err
			}
		}
		return nil
	})
}
