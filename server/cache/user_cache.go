package cache

import (
	"errors"
	"time"

	"bbs-go/model"
	"bbs-go/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
)

type userCache struct {
	cache          cache.LoadingCache
	scoreRankCache cache.LoadingCache
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserRepository.Get(simple.DB(), key2Int64(key))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		scoreRankCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserRepository.Find(simple.DB(), simple.NewSqlCnd().Desc("score").Limit(10))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(10*time.Minute),
		),
	}
}

func (c *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	val, err := c.cache.Get(userId)
	if err != nil {
		return nil
	}
	return val.(*model.User)
}

func (c *userCache) Invalidate(userId int64) {
	c.cache.Invalidate(userId)
}

func (c *userCache) GetScoreRank() []model.User {
	val, err := c.scoreRankCache.Get("data")
	if err != nil {
		return nil
	}
	return val.([]model.User)
}
