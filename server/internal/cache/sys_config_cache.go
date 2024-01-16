package cache

import (
	"errors"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

type sysConfigCache struct {
	cache cache.LoadingCache
}

var SysConfigCache = newSysConfigCache()

func newSysConfigCache() *sysConfigCache {
	return &sysConfigCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.SysConfigRepository.GetByKey(sqls.DB(), key.(string))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *sysConfigCache) Get(key string) *models.SysConfig {
	val, err := c.cache.Get(key)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.(*models.SysConfig)
	}
	return nil
}

func (c *sysConfigCache) GetValue(key string) string {
	sysConfig := c.Get(key)
	if sysConfig == nil {
		return ""
	}
	return sysConfig.Value
}

func (c *sysConfigCache) Invalidate(key string) {
	c.cache.Invalidate(key)
}
