package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/arrs"
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
	if max := s.FindOne(sqls.NewCnd().Eq("parent_id", parentId).Eq("status", constants.StatusOk).Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *menuService) GetUserMenuIds(userId int64) []int64 {
	roleIds := UserRoleService.GetUserRoleIds(userId)
	if len(roleIds) == 0 {
		return nil
	}
	return RoleMenuService.GetMenuIdsByRoles(roleIds)
}

func (s *menuService) GetUserMenus(user *models.User) (ret []models.Menu) {
	menuIds := s.GetUserMenuIds(user.Id)
	if len(menuIds) == 0 {
		return nil
	}

	menus := s.Find(sqls.NewCnd().Eq("status", constants.StatusOk).Asc("sort_no").Desc("id"))
	menusMap := make(map[int64]models.Menu, len(menus))
	for _, menu := range menus {
		menusMap[menu.Id] = menu
	}

	var showMenuIds []int64
	for _, menuId := range menuIds {
		menuPath := s.GetMenuPath(menuId, menusMap)
		showMenuIds = append(showMenuIds, menuPath...)
	}

	for _, menu := range menus {
		if arrs.Contains(showMenuIds, menu.Id) {
			ret = append(ret, menu)
		}
	}
	return ret
}

// GetMenuPath 获取菜单路径
func (s *menuService) GetMenuPath(menuId int64, menusMap map[int64]models.Menu) (ret []int64) {
	if menuId <= 0 {
		return
	}
	for {
		menu, found := menusMap[menuId]
		if !found {
			break
		}
		ret = append(ret, menu.Id)
		if menu.ParentId > 0 {
			menuId = menu.ParentId
		} else {
			break
		}
	}
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return
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
