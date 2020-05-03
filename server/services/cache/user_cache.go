package cache

import (
	"time"

	"bbs-go/model"
	"bbs-go/repositories"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
)

type userCache struct {
	cache      cache.LoadingCache
	scoreCache cache.LoadingCache
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserRepository.Get(simple.DB(), key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
		scoreCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, err error) {
				userScore := repositories.UserScoreRepository.FindOne(simple.DB(),
					simple.NewSqlCnd().Eq("user_id", key2Int64(key)))
				if userScore == nil {
					value = 0
				} else {
					value = userScore.Score
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
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

func (c *userCache) GetScore(userId int64) int {
	val, err := c.scoreCache.Get(userId)
	if err != nil {
		return 0
	}
	return val.(int)
}

func (c *userCache) InvalidateScore(userId int64) {
	c.scoreCache.Invalidate(userId)
}
