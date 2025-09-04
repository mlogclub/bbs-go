package cache

import (
	"context"
	"log/slog"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisForbiddenWordCache struct {
	cache redis.Cache
}

var ForbiddenWordCacheRedis = newRedisForbiddenWordCache()

func newRedisForbiddenWordCache() *redisForbiddenWordCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisForbiddenWordCache{
		cache: cache,
	}
}

func (c *redisForbiddenWordCache) Get() []models.ForbiddenWord {
	return c.GetAll()
}

func (c *redisForbiddenWordCache) GetAll() []models.ForbiddenWord {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getAllFromDB()
	}
	
	ctx := context.Background()
	
	var words []models.ForbiddenWord
	err := c.cache.GetJSON(ctx, redis.ForbiddenWordsKey, &words)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		words = c.getAllFromDB()
		if len(words) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), redis.ForbiddenWordsKey, words, redis.ForbiddenWordsExpire)
			}()
		}
	}
	
	return words
}

func (c *redisForbiddenWordCache) getAllFromDB() []models.ForbiddenWord {
	if c == nil {
		return repositories.ForbiddenWordRepository.Find(sqls.DB(), sqls.NewCnd())
	}
	return repositories.ForbiddenWordRepository.Find(sqls.DB(), sqls.NewCnd())
}

func (c *redisForbiddenWordCache) Refresh() {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	words := c.getAllFromDB()
	
	err := c.cache.SetJSON(ctx, redis.ForbiddenWordsKey, words, redis.ForbiddenWordsExpire)
	if err != nil {
		slog.Warn("Redis refresh forbidden words cache failed", "error", err)
	}
}

func (c *redisForbiddenWordCache) Invalidate() {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	err := c.cache.Del(ctx, redis.ForbiddenWordsKey)
	if err != nil {
		slog.Warn("Redis delete forbidden words cache failed", "error", err)
	}
}

// 批量预热违禁词缓存
func (c *redisForbiddenWordCache) WarmUp() {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	words := c.getAllFromDB()
	if len(words) == 0 {
		return
	}
	
	ctx := context.Background()
	err := c.cache.SetJSON(ctx, redis.ForbiddenWordsKey, words, redis.ForbiddenWordsExpire)
	if err != nil {
		slog.Warn("Redis warm up forbidden words cache failed", "error", err)
	} else {
		slog.Info("Forbidden words cache warmed up", "count", len(words))
	}
}