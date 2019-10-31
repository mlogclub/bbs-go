package cache

import (
	"github.com/goburrow/cache"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/simple"
	"time"
)

type sysConfigCache struct {
	cache cache.LoadingCache
}

var SysConfigCache = newSysConfigCache()

func newSysConfigCache() *sysConfigCache {
	return &sysConfigCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.SysConfigRepository.GetByKey(simple.DB(), key.(string))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (this *sysConfigCache) Get(key string) *model.SysConfig {
	val, err := this.cache.Get(key)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*model.SysConfig)
	}
	return nil
}

func (this *sysConfigCache) GetValue(key string) string {
	sysConfig := this.Get(key)
	if sysConfig == nil {
		return ""
	}
	return sysConfig.Value
}

func (this *sysConfigCache) Invalidate(key string) {
	this.cache.Invalidate(key)
}
