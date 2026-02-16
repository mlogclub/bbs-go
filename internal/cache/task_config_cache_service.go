package cache

import (
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
)

const taskConfigCacheKey = "all_task_configs"

// taskConfigCacheService keeps TaskConfig in memory for quick reads.
type taskConfigCacheService struct {
	cache cache.LoadingCache
}

var TaskConfigCacheService = newTaskConfigCacheService()

func newTaskConfigCacheService() *taskConfigCacheService {
	return &taskConfigCacheService{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.TaskConfigRepository.Find(sqls.DB(), sqls.NewCnd().
					Eq("status", constants.StatusOk).
					Asc("sort_no"))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

// GetAll returns a defensive copy of cached TaskConfig list.
func (s *taskConfigCacheService) GetAll() []models.TaskConfig {
	val, err := s.cache.Get(taskConfigCacheKey)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.TaskConfig)
	cp := make([]models.TaskConfig, len(list))
	copy(cp, list)
	return cp
}

// GetByEventType returns cached TaskConfigs filtered by eventType.
func (s *taskConfigCacheService) GetByEventType(eventType string) []models.TaskConfig {
	if strs.IsBlank(eventType) {
		return nil
	}
	all := s.GetAll()
	if len(all) == 0 {
		return nil
	}
	ret := make([]models.TaskConfig, 0, len(all))
	for i := range all {
		if all[i].EventType == eventType {
			ret = append(ret, all[i])
		}
	}
	return ret
}

// Reload refreshes the cache immediately.
func (s *taskConfigCacheService) Reload() {
	s.cache.Refresh(taskConfigCacheKey)
}
