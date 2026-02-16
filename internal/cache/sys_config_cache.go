package cache

import (
	"errors"
	"log/slog"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"

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
				if ret := repositories.SysConfigRepository.GetByKey(sqls.DB(), key.(string)); ret != nil {
					value = ret
				} else {
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

func (c *sysConfigCache) GetStr(key string) string {
	if t := c.Get(key); t != nil {
		return t.Value
	}
	return ""
}

func (c *sysConfigCache) GetBool(key string) bool {
	str := c.GetStr(key)
	return cast.ToBool(str)
}

func (c *sysConfigCache) GetInt(key string) int {
	str := c.GetStr(key)
	return cast.ToInt(str)
}

func (c *sysConfigCache) GetInt64(key string) int64 {
	str := c.GetStr(key)
	return cast.ToInt64(str)
}

func (c *sysConfigCache) GetStrArr(key string) (ret []string) {
	str := c.GetStr(key)
	if err := jsons.Parse(str, &ret); err != nil {
		slog.Warn("config value error", slog.Any("key", key), slog.Any("err", err))
	}
	return
}

func (c *sysConfigCache) Invalidate(key string) {
	c.cache.Invalidate(key)
}
