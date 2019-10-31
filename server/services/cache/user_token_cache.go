package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var UserTokenCache = newUserTokenCache()

type userTokenCache struct {
	cache cache.LoadingCache
}

func newUserTokenCache() *userTokenCache {
	return &userTokenCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.UserTokenRepository.GetByToken(simple.DB(), key.(string))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(60*time.Minute),
		),
	}
}

func (this *userTokenCache) Get(token string) *model.UserToken {
	if len(token) == 0 {
		return nil
	}
	val, err := this.cache.Get(token)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.UserToken)
	}
	return nil
}

func (this *userTokenCache) Invalidate(token string) {
	this.cache.Invalidate(token)
}
