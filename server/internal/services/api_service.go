package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/repositories"
	"log/slog"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var ApiService = newApiService()

func newApiService() *apiService {
	return &apiService{}
}

type apiService struct {
}

func (s *apiService) Get(id int64) *models.Api {
	return repositories.ApiRepository.Get(sqls.DB(), id)
}

func (s *apiService) Take(where ...interface{}) *models.Api {
	return repositories.ApiRepository.Take(sqls.DB(), where...)
}

func (s *apiService) Find(cnd *sqls.Cnd) []models.Api {
	return repositories.ApiRepository.Find(sqls.DB(), cnd)
}

func (s *apiService) FindOne(cnd *sqls.Cnd) *models.Api {
	return repositories.ApiRepository.FindOne(sqls.DB(), cnd)
}

func (s *apiService) FindPageByParams(params *params.QueryParams) (list []models.Api, paging *sqls.Paging) {
	return repositories.ApiRepository.FindPageByParams(sqls.DB(), params)
}

func (s *apiService) FindPageByCnd(cnd *sqls.Cnd) (list []models.Api, paging *sqls.Paging) {
	return repositories.ApiRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *apiService) Count(cnd *sqls.Cnd) int64 {
	return repositories.ApiRepository.Count(sqls.DB(), cnd)
}

func (s *apiService) Create(t *models.Api) error {
	return repositories.ApiRepository.Create(sqls.DB(), t)
}

func (s *apiService) Update(t *models.Api) error {
	return repositories.ApiRepository.Update(sqls.DB(), t)
}

func (s *apiService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.ApiRepository.Updates(sqls.DB(), id, columns)
}

func (s *apiService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.ApiRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *apiService) Delete(id int64) {
	repositories.ApiRepository.Delete(sqls.DB(), id)
}

func (s *apiService) GetByPath(path string) *models.Api {
	return s.FindOne(sqls.NewCnd().Where("path = ?", path))
}

func (s *apiService) GetByName(name string) *models.Api {
	return s.FindOne(sqls.NewCnd().Where("name = ?", name))
}

func (s *apiService) Init(list []dto.ApiRoute) {
	now := dates.NowTimestamp()
	for _, item := range list {
		if s.GetByPath(item.Path) != nil {
			continue
		}

		if err := s.Create(&models.Api{
			Name:       item.Name,
			Method:     item.Method,
			Path:       item.Path,
			CreateTime: now,
			UpdateTime: now,
		}); err != nil {
			slog.Error("Create api error", slog.Any("error", err))
		} else {
			slog.Info("Create api: " + item.Method + " " + item.Path + " " + item.Name)
		}
	}
}

func (s *apiService) GetByMenuId(menuId int64) (list []models.Api) {
	err := sqls.DB().Raw(`
	select a.* from t_api a left join t_menu_api ma on a.id = ma.api_id where ma.menu_id = ? order by a.id
	`, menuId).Scan(&list).Error
	if err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
	}
	return
}

func (s *apiService) GetByUserId(userId int64) (list []models.Api) {
	err := sqls.DB().Raw(`
	select api.* from t_api api
	left join t_menu_api menu_api on api.id = menu_api.api_id
	left join t_role_menu role_menu on role_menu.menu_id = menu_api.menu_id
	left join t_user_role user_role on user_role.role_id = role_menu.role_id
	where user_role.user_id = ? order by api.id
	`, userId).Scan(&list).Error
	if err != nil {
		slog.Error(err.Error(), slog.Any("error", err))
	}
	return
}
