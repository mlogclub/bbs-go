package services

import (
	"errors"

	"bbs-go/internal/cache"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
)

var LevelConfigService = newLevelConfigService()

func newLevelConfigService() *levelConfigService {
	return &levelConfigService{}
}

type levelConfigService struct {
}

func (s *levelConfigService) Get(id int64) *models.LevelConfig {
	return repositories.LevelConfigRepository.Get(sqls.DB(), id)
}

func (s *levelConfigService) Take(where ...interface{}) *models.LevelConfig {
	return repositories.LevelConfigRepository.Take(sqls.DB(), where...)
}

func (s *levelConfigService) Find(cnd *sqls.Cnd) []models.LevelConfig {
	return repositories.LevelConfigRepository.Find(sqls.DB(), cnd)
}

func (s *levelConfigService) FindOne(cnd *sqls.Cnd) *models.LevelConfig {
	return repositories.LevelConfigRepository.FindOne(sqls.DB(), cnd)
}

func (s *levelConfigService) FindPageByParams(params *params.QueryParams) (list []models.LevelConfig, paging *sqls.Paging) {
	return repositories.LevelConfigRepository.FindPageByParams(sqls.DB(), params)
}

func (s *levelConfigService) FindPageByCnd(cnd *sqls.Cnd) (list []models.LevelConfig, paging *sqls.Paging) {
	return repositories.LevelConfigRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *levelConfigService) Count(cnd *sqls.Cnd) int64 {
	return repositories.LevelConfigRepository.Count(sqls.DB(), cnd)
}

// SaveAll 批量保存等级配置（Level 必须从 1 开始且连续，NeedExp 必须严格递增）
func (s *levelConfigService) SaveAll(items []models.LevelConfig) error {
	if err := validateLevelConfigItems(items); err != nil {
		return err
	}

	now := dates.NowTimestamp()
	err := sqls.WithTransaction(func(ctx *sqls.TxContext) error {
		existing := repositories.LevelConfigRepository.Find(ctx.Tx, sqls.NewCnd().Asc("level"))
		existingByLevel := make(map[int]int64, len(existing))
		maxExistingLevel := 0
		for _, e := range existing {
			existingByLevel[e.Level] = e.Id
			if e.Level > maxExistingLevel {
				maxExistingLevel = e.Level
			}
		}

		maxLevel := len(items)
		for i := range items {
			items[i].Status = constants.StatusOk
			items[i].UpdateTime = now
			if items[i].CreateTime == 0 {
				items[i].CreateTime = now
			}

			if id := existingByLevel[items[i].Level]; id > 0 {
				items[i].Id = id
				if err := repositories.LevelConfigRepository.Update(ctx.Tx, &items[i]); err != nil {
					return err
				}
			} else {
				items[i].Id = 0
				if err := repositories.LevelConfigRepository.Create(ctx.Tx, &items[i]); err != nil {
					return err
				}
			}
		}

		// 软删除超出范围的历史等级（只允许删尾部）
		if maxExistingLevel > maxLevel {
			for level := maxLevel + 1; level <= maxExistingLevel; level++ {
				if id := existingByLevel[level]; id > 0 {
					if err := repositories.LevelConfigRepository.Updates(ctx.Tx, id, map[string]interface{}{
						"status":      constants.StatusDeleted,
						"update_time": now,
					}); err != nil {
						return err
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	cache.LevelConfigCache.Reload()
	return nil
}

// GetAllFromCache returns cached level configs (ordered by level).
func (s *levelConfigService) GetAllFromCache() []models.LevelConfig {
	return cache.LevelConfigCache.GetAll()
}

// GetByLevelFromCache returns the cached config for the given level.
func (s *levelConfigService) GetByLevelFromCache(level int) *models.LevelConfig {
	return cache.LevelConfigCache.GetByLevel(level)
}

func validateLevelConfigItems(items []models.LevelConfig) error {
	if len(items) == 0 {
		return errors.New("level config is empty")
	}
	for i := range items {
		expectedLevel := i + 1
		if items[i].Level != expectedLevel {
			return errors.New("level must start from 1 and be continuous")
		}
		if i == 0 && items[i].NeedExp != 0 {
			return errors.New("level 1 needExp must be 0")
		}
		if items[i].NeedExp < 0 {
			return errors.New("needExp must be >= 0")
		}
		if i > 0 && items[i].NeedExp <= items[i-1].NeedExp {
			return errors.New("needExp must be strictly increasing")
		}
	}
	return nil
}
