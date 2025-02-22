package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var MenuApiService = newMenuApiService()

func newMenuApiService() *menuApiService {
	return &menuApiService{}
}

type menuApiService struct {
}

func (s *menuApiService) Get(id int64) *models.MenuApi {
	return repositories.MenuApiRepository.Get(sqls.DB(), id)
}

func (s *menuApiService) Take(where ...interface{}) *models.MenuApi {
	return repositories.MenuApiRepository.Take(sqls.DB(), where...)
}

func (s *menuApiService) Find(cnd *sqls.Cnd) []models.MenuApi {
	return repositories.MenuApiRepository.Find(sqls.DB(), cnd)
}

func (s *menuApiService) FindOne(cnd *sqls.Cnd) *models.MenuApi {
	return repositories.MenuApiRepository.FindOne(sqls.DB(), cnd)
}

func (s *menuApiService) FindPageByParams(params *params.QueryParams) (list []models.MenuApi, paging *sqls.Paging) {
	return repositories.MenuApiRepository.FindPageByParams(sqls.DB(), params)
}

func (s *menuApiService) FindPageByCnd(cnd *sqls.Cnd) (list []models.MenuApi, paging *sqls.Paging) {
	return repositories.MenuApiRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *menuApiService) Count(cnd *sqls.Cnd) int64 {
	return repositories.MenuApiRepository.Count(sqls.DB(), cnd)
}

func (s *menuApiService) Create(t *models.MenuApi) error {
	return repositories.MenuApiRepository.Create(sqls.DB(), t)
}

func (s *menuApiService) Update(t *models.MenuApi) error {
	return repositories.MenuApiRepository.Update(sqls.DB(), t)
}

func (s *menuApiService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.MenuApiRepository.Updates(sqls.DB(), id, columns)
}

func (s *menuApiService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.MenuApiRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *menuApiService) Delete(id int64) {
	repositories.MenuApiRepository.Delete(sqls.DB(), id)
}

func (s *menuApiService) SetMenuApis(menuId int64, apiIds []int64) {
	sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.MenuApi{}, "menu_id = ?", menuId).Error; err != nil {
			return err
		}
		for _, apiId := range apiIds {
			if err := repositories.MenuApiRepository.Create(tx, &models.MenuApi{
				MenuId:     menuId,
				ApiId:      apiId,
				CreateTime: dates.NowTimestamp(),
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *menuApiService) SetMenuApis2(menuId int64, apis []models.Api) {
	var apiIds []int64
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}
	s.SetMenuApis(menuId, apiIds)
}
