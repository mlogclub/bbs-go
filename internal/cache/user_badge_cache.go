package cache

import (
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"
)

// userBadgeCache wraps a LoadingCache: key = userId, value = []models.UserBadge.
type userBadgeCache struct {
	cache cache.LoadingCache
}

var UserBadgeCache = newUserBadgeCache()

func newUserBadgeCache() *userBadgeCache {
	return &userBadgeCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				uid := key2Int64(key)
				value = repositories.UserBadgeRepository.Find(sqls.DB(),
					sqls.NewCnd().Eq("user_id", uid))
				return
			},
			cache.WithMaximumSize(10000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

// GetByUser returns cached UserBadge list for the user (copy to avoid mutation).
func (c *userBadgeCache) GetByUser(userId int64) []models.UserBadge {
	if userId <= 0 {
		return nil
	}
	val, err := c.cache.Get(userId)
	if err != nil || val == nil {
		return nil
	}
	list := val.([]models.UserBadge)
	cp := make([]models.UserBadge, len(list))
	copy(cp, list)
	return cp
}

// Invalidate clears the cache for the given user (next GetByUser will reload from DB).
func (c *userBadgeCache) Invalidate(userId int64) {
	if userId <= 0 {
		return
	}
	c.cache.Invalidate(userId)
}
