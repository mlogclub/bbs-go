package services

import (
	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

var TaskConfigService = newTaskConfigService()

func newTaskConfigService() *taskConfigService {
	return &taskConfigService{}
}

type taskConfigService struct {
}

func (s *taskConfigService) Get(id int64) *models.TaskConfig {
	return repositories.TaskConfigRepository.Get(sqls.DB(), id)
}

func (s *taskConfigService) Take(where ...interface{}) *models.TaskConfig {
	return repositories.TaskConfigRepository.Take(sqls.DB(), where...)
}

func (s *taskConfigService) Find(cnd *sqls.Cnd) []models.TaskConfig {
	return repositories.TaskConfigRepository.Find(sqls.DB(), cnd)
}

func (s *taskConfigService) FindOne(cnd *sqls.Cnd) *models.TaskConfig {
	return repositories.TaskConfigRepository.FindOne(sqls.DB(), cnd)
}

func (s *taskConfigService) FindPageByParams(params *params.QueryParams) (list []models.TaskConfig, paging *sqls.Paging) {
	return repositories.TaskConfigRepository.FindPageByParams(sqls.DB(), params)
}

func (s *taskConfigService) FindPageByCnd(cnd *sqls.Cnd) (list []models.TaskConfig, paging *sqls.Paging) {
	return repositories.TaskConfigRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *taskConfigService) Count(cnd *sqls.Cnd) int64 {
	return repositories.TaskConfigRepository.Count(sqls.DB(), cnd)
}

func (s *taskConfigService) Create(t *models.TaskConfig) error {
	if err := repositories.TaskConfigRepository.Create(sqls.DB(), t); err != nil {
		return err
	}

	cache.TaskConfigCacheService.Reload()
	return nil
}

func (s *taskConfigService) Update(t *models.TaskConfig) error {
	if err := repositories.TaskConfigRepository.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.TaskConfigCacheService.Reload()
	return nil
}

func (s *taskConfigService) Updates(id int64, columns map[string]interface{}) error {
	if err := repositories.TaskConfigRepository.Updates(sqls.DB(), id, columns); err != nil {
		return err
	}
	cache.TaskConfigCacheService.Reload()
	return nil
}

func (s *taskConfigService) GetNextSortNo() int {
	if max := s.FindOne(sqls.NewCnd().Eq("status", constants.StatusOk).Desc("sort_no")); max != nil {
		return max.SortNo + 1
	}
	return 0
}

func (s *taskConfigService) UpdateSort(ids []int64) error {
	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := repositories.TaskConfigRepository.UpdateColumn(tx, id, "sort_no", i); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	cache.TaskConfigCacheService.Reload()
	return nil
}
