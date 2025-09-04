package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/spf13/cast"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/common/jsons"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisSysConfigCache struct {
	cache redis.Cache
}

var SysConfigCacheRedis = newRedisSysConfigCache()

func newRedisSysConfigCache() *redisSysConfigCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisSysConfigCache{
		cache: cache,
	}
}

func (c *redisSysConfigCache) Get(key string) *models.SysConfig {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(key)
	}
	
	ctx := context.Background()
	cacheKey := fmt.Sprintf(redis.SysConfigCacheKey, key)
	
	// 使用防缓存击穿获取数据
	result, err := c.cache.GetWithBreakdown(ctx, cacheKey, redis.ConfigCacheExpire, func() (interface{}, error) {
		return c.getFromDB(key), nil
	})
	
	if err != nil {
		slog.Warn("Redis get sys config failed, fallback to database", "key", key, "error", err)
		return c.getFromDB(key)
	}
	
	if result == nil {
		return nil
	}
	
	config, ok := result.(*models.SysConfig)
	if !ok {
		// 类型断言失败，从数据库重新加载
		return c.getFromDB(key)
	}
	
	return config
}

func (c *redisSysConfigCache) getFromDB(key string) *models.SysConfig {
	if c == nil {
		return repositories.SysConfigRepository.GetByKey(sqls.DB(), key)
	}
	return repositories.SysConfigRepository.GetByKey(sqls.DB(), key)
}

func (c *redisSysConfigCache) GetStr(key string) string {
	if t := c.Get(key); t != nil {
		return t.Value
	}
	return ""
}

func (c *redisSysConfigCache) GetBool(key string) bool {
	str := c.GetStr(key)
	return cast.ToBool(str)
}

func (c *redisSysConfigCache) GetInt(key string) int {
	str := c.GetStr(key)
	return cast.ToInt(str)
}

func (c *redisSysConfigCache) GetInt64(key string) int64 {
	str := c.GetStr(key)
	return cast.ToInt64(str)
}

func (c *redisSysConfigCache) GetStrArr(key string) (ret []string) {
	str := c.GetStr(key)
	if err := jsons.Parse(str, &ret); err != nil {
		slog.Warn("config value error", slog.Any("key", key), slog.Any("err", err))
	}
	return
}

func (c *redisSysConfigCache) Invalidate(key string) {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	cacheKey := fmt.Sprintf(redis.SysConfigCacheKey, key)
	
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		slog.Warn("Redis delete sys config cache failed", "key", key, "error", err)
	}
}

func (c *redisSysConfigCache) Set(key string, config *models.SysConfig) {
	if c == nil || c.cache == nil || config == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	cacheKey := fmt.Sprintf(redis.SysConfigCacheKey, key)
	
	err := c.cache.SetJSON(ctx, cacheKey, config, redis.ConfigCacheExpire)
	if err != nil {
		slog.Warn("Redis set sys config cache failed", "key", key, "error", err)
	}
}

// 批量预热系统配置缓存
func (c *redisSysConfigCache) WarmUp() {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	configs := repositories.SysConfigRepository.Find(sqls.DB(), sqls.NewCnd())
	
	if len(configs) == 0 {
		return
	}
	
	pipe := c.cache.Pipeline()
	for _, config := range configs {
		key := fmt.Sprintf(redis.SysConfigCacheKey, config.Key)
		// 使用JSON序列化设置缓存
		if data, err := json.Marshal(&config); err == nil {
			pipe.Set(ctx, key, data, redis.ConfigCacheExpire)
		}
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Warn("Redis warm up sys config cache failed", "error", err)
	} else {
		slog.Info("Sys config cache warmed up", "count", len(configs))
	}
}