package cache

import (
	"context"
	"fmt"
	"log/slog"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisUserTokenCache struct {
	cache redis.Cache
}

var UserTokenCacheRedis = newRedisUserTokenCache()

func newRedisUserTokenCache() *redisUserTokenCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisUserTokenCache{
		cache: cache,
	}
}

func (c *redisUserTokenCache) Get(token string) *models.UserToken {
	if len(token) == 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(token)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.UserTokenCacheKey, token)
	
	// 使用防缓存击穿获取数据
	result, err := c.cache.GetWithBreakdown(ctx, key, redis.TokenCacheExpire, func() (interface{}, error) {
		return c.getFromDB(token), nil
	})
	
	if err != nil {
		slog.Warn("Redis get user token failed, fallback to database", "token", token, "error", err)
		return c.getFromDB(token)
	}
	
	if result == nil {
		return nil
	}
	
	userToken, ok := result.(*models.UserToken)
	if !ok {
		// 类型断言失败，从数据库重新加载
		return c.getFromDB(token)
	}
	
	return userToken
}

func (c *redisUserTokenCache) getFromDB(token string) *models.UserToken {
	if c == nil {
		return repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
	}
	return repositories.UserTokenRepository.GetByToken(sqls.DB(), token)
}

func (c *redisUserTokenCache) Invalidate(token string) {
	if c == nil || c.cache == nil || len(token) == 0 || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.UserTokenCacheKey, token)
	
	err := c.cache.Del(ctx, key)
	if err != nil {
		slog.Warn("Redis delete user token cache failed", "token", token, "error", err)
	}
}

func (c *redisUserTokenCache) Set(token string, userToken *models.UserToken) {
	if c == nil || c.cache == nil || len(token) == 0 || userToken == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.UserTokenCacheKey, token)
	
	err := c.cache.SetJSON(ctx, key, userToken, redis.TokenCacheExpire)
	if err != nil {
		slog.Warn("Redis set user token cache failed", "token", token, "error", err)
	}
}