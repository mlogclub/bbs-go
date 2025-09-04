package cache

import (
	"context"
	"log/slog"
	"time"

	"bbs-go/internal/pkg/redis"
)

// EnhancedCache 增强型缓存接口，支持降级和错误处理
type EnhancedCache struct {
	redis      redis.Cache
	name       string
	fallbackFn func(key string) (interface{}, error)
}

// NewEnhancedCache 创建增强型缓存
func NewEnhancedCache(name string, fallbackFn func(key string) (interface{}, error)) *EnhancedCache {
	return &EnhancedCache{
		redis:      redis.NewRedisCache(),
		name:       name,
		fallbackFn: fallbackFn,
	}
}

// GetWithFallback 获取缓存数据，支持降级
func (ec *EnhancedCache) GetWithFallback(ctx context.Context, key string, expire time.Duration) (interface{}, error) {
	if ec.redis == nil {
		slog.Warn("Redis cache not available, using fallback", "cache", ec.name, "key", key)
		return ec.fallbackFn(key)
	}

	// 尝试从缓存获取
	var result interface{}
	err := ec.redis.GetJSON(ctx, key, &result)
	if err == nil {
		return result, nil
	}

	// 缓存未命中或错误，记录日志并使用降级
	if err == redis.ErrCacheNotFound {
		slog.Debug("Cache miss, using fallback", "cache", ec.name, "key", key)
	} else if err == redis.ErrCircuitOpen {
		slog.Warn("Circuit breaker open, using fallback", "cache", ec.name, "key", key)
	} else {
		slog.Warn("Cache error, using fallback", "cache", ec.name, "key", key, "error", err)
	}

	// 使用降级函数获取数据
	data, fallbackErr := ec.fallbackFn(key)
	if fallbackErr != nil {
		return nil, fallbackErr
	}

	// 异步设置缓存
	if data != nil {
		go ec.asyncSetCache(key, data, expire)
	}

	return data, nil
}

// SafeSet 安全设置缓存，忽略错误
func (ec *EnhancedCache) SafeSet(ctx context.Context, key string, value interface{}, expire time.Duration) {
	if ec.redis == nil {
		return
	}

	if err := ec.redis.SetJSON(ctx, key, value, expire); err != nil {
		slog.Warn("Failed to set cache", "cache", ec.name, "key", key, "error", err)
	}
}

// SafeDelete 安全删除缓存，忽略错误
func (ec *EnhancedCache) SafeDelete(ctx context.Context, keys ...string) {
	if ec.redis == nil {
		return
	}

	if err := ec.redis.Del(ctx, keys...); err != nil {
		slog.Warn("Failed to delete cache", "cache", ec.name, "keys", keys, "error", err)
	}
}

// asyncSetCache 异步设置缓存
func (ec *EnhancedCache) asyncSetCache(key string, value interface{}, expire time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ec.redis.SetJSON(ctx, key, value, expire); err != nil {
		slog.Warn("Failed to set cache asynchronously", "cache", ec.name, "key", key, "error", err)
	}
}

// IsHealthy 检查缓存是否健康
func (ec *EnhancedCache) IsHealthy() bool {
	if ec.redis == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 尝试ping Redis
	_, err := ec.redis.Get(ctx, "health_check")
	return err == nil || err == redis.ErrCacheNotFound
}

// WarmUp 预热缓存
func (ec *EnhancedCache) WarmUp(keys []string, expire time.Duration) {
	if ec.redis == nil || ec.fallbackFn == nil {
		return
	}

	slog.Info("Starting cache warm up", "cache", ec.name, "keys_count", len(keys))

	for _, key := range keys {
		go func(k string) {
			data, err := ec.fallbackFn(k)
			if err != nil {
				slog.Warn("Failed to warm up cache", "cache", ec.name, "key", k, "error", err)
				return
			}

			ec.SafeSet(context.Background(), k, data, expire)
		}(key)
	}
}