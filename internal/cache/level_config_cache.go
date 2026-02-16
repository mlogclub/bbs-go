package cache

import (
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"
)

const levelConfigCacheKey = "all_level_configs"

// levelConfigCache wraps a LoadingCache to keep all level configs in memory.
type levelConfigCache struct {
	cache cache.LoadingCache
}

var LevelConfigCache = newLevelConfigCache()

func newLevelConfigCache() *levelConfigCache {
	return &levelConfigCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.LevelConfigRepository.Find(sqls.DB(), sqls.NewCnd().Asc("level"))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

// GetAll returns cached LevelConfig list (copied to avoid external mutation).
func (c *levelConfigCache) GetAll() []models.LevelConfig {
	val, err := c.cache.Get(levelConfigCacheKey)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.LevelConfig)
	cp := make([]models.LevelConfig, len(list))
	copy(cp, list)
	return cp
}

// GetByLevel returns the LevelConfig for a given level from cache.
func (c *levelConfigCache) GetByLevel(level int) *models.LevelConfig {
	val, err := c.cache.Get(levelConfigCacheKey)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.LevelConfig)
	for i := range list {
		if list[i].Level == level {
			// return a copy to avoid mutation
			ret := list[i]
			return &ret
		}
	}
	return nil
}

// Reload refreshes the cache immediately.
func (c *levelConfigCache) Reload() {
	c.cache.Refresh(levelConfigCacheKey)
}
