package cache

import (
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"
)

const badgeCacheKey = "all_badges"

// badgeCache wraps a LoadingCache to keep all badges (status=ok) in memory.
type badgeCache struct {
	cache cache.LoadingCache
}

var BadgeCache = newBadgeCache()

func newBadgeCache() *badgeCache {
	return &badgeCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.BadgeRepository.Find(sqls.DB(), sqls.NewCnd().
					Eq("status", constants.StatusOk).
					Asc("sort_no").
					Desc("id"))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithExpireAfterAccess(180*time.Minute),
		),
	}
}

// GetAll returns cached Badge list (copied to avoid external mutation).
func (c *badgeCache) GetAll() []models.Badge {
	val, err := c.cache.Get(badgeCacheKey)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.Badge)
	cp := make([]models.Badge, len(list))
	copy(cp, list)
	return cp
}

// GetByID returns the Badge for a given id from cache.
func (c *badgeCache) GetByID(id int64) *models.Badge {
	val, err := c.cache.Get(badgeCacheKey)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.Badge)
	for i := range list {
		if list[i].Id == id {
			ret := list[i]
			return &ret
		}
	}
	return nil
}

// Reload refreshes the cache immediately.
func (c *badgeCache) Reload() {
	c.cache.Refresh(badgeCacheKey)
}
