package cache

import (
	"time"

	"bbs-go/model"
	"bbs-go/repositories"

	"github.com/goburrow/cache"
	"bbs-go/simple"
)

type userCache struct {
	cache cache.LoadingCache
}

var UserCache = newUserCache()

func newUserCache() *userCache {
	return &userCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserRepository.Get(simple.DB(), Key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (this *userCache) Get(userId int64) *model.User {
	if userId <= 0 {
		return nil
	}
	val, err := this.cache.Get(userId)
	if err != nil {
		return nil
	}
	return val.(*model.User)
}

func (this *userCache) Invalidate(userId int64) {
	this.cache.Invalidate(userId)
}
