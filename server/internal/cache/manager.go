package cache

import (
	"context"
	"log/slog"
	"os"

	"bbs-go/internal/pkg/redis"
)

// CacheManager 缓存管理器
type CacheManager struct {
	redisEnabled bool
}

var Manager = &CacheManager{}

// Init 初始化缓存管理器
func (m *CacheManager) Init() {
	m.redisEnabled = redis.IsRedisEnabled()
	
	if !m.redisEnabled {
		slog.Error("Redis is required but not enabled. Please configure Redis properly.")
		os.Exit(1)
	}
	
	slog.Info("Cache manager initialized with Redis backend")
	// 预热缓存
	m.WarmUp()
}

// IsRedisEnabled 检查Redis是否启用
func (m *CacheManager) IsRedisEnabled() bool {
	return m.redisEnabled
}

// GetUserCache 获取用户缓存实例
func (m *CacheManager) GetUserCache() interface{} {
	if UserCacheRedis == nil {
		slog.Error("UserCacheRedis is not initialized")
		os.Exit(1)
	}
	return UserCacheRedis
}

// GetUserTokenCache 获取用户Token缓存实例
func (m *CacheManager) GetUserTokenCache() interface{} {
	if UserTokenCacheRedis == nil {
		slog.Error("UserTokenCacheRedis is not initialized")
		os.Exit(1)
	}
	return UserTokenCacheRedis
}

// GetSysConfigCache 获取系统配置缓存实例
func (m *CacheManager) GetSysConfigCache() interface{} {
	if SysConfigCacheRedis == nil {
		slog.Error("SysConfigCacheRedis is not initialized")
		os.Exit(1)
	}
	return SysConfigCacheRedis
}

// GetTopicCache 获取话题缓存实例
func (m *CacheManager) GetTopicCache() interface{} {
	if TopicCacheRedis == nil {
		slog.Error("TopicCacheRedis is not initialized")
		os.Exit(1)
	}
	return TopicCacheRedis
}

// GetTagCache 获取标签缓存实例
func (m *CacheManager) GetTagCache() interface{} {
	if TagCacheRedis == nil {
		slog.Error("TagCacheRedis is not initialized")
		os.Exit(1)
	}
	return TagCacheRedis
}

// GetForbiddenWordCache 获取违禁词缓存实例
func (m *CacheManager) GetForbiddenWordCache() interface{} {
	if ForbiddenWordCacheRedis == nil {
		slog.Error("ForbiddenWordCacheRedis is not initialized")
		os.Exit(1)
	}
	return ForbiddenWordCacheRedis
}

// GetArticleTagCache 获取文章标签缓存实例
func (m *CacheManager) GetArticleTagCache() interface{} {
	if ArticleTagCacheRedis == nil {
		slog.Error("ArticleTagCacheRedis is not initialized")
		os.Exit(1)
	}
	return ArticleTagCacheRedis
}

// WarmUp 预热缓存
func (m *CacheManager) WarmUp() {
	if !m.redisEnabled {
		return
	}
	
	slog.Info("Starting cache warm up...")
	
	// 预热系统配置
	if SysConfigCacheRedis != nil {
		go SysConfigCacheRedis.WarmUp()
	}
	
	// 预热标签缓存
	if TagCacheRedis != nil {
		go TagCacheRedis.WarmUp()
	}
	
	// 预热违禁词缓存
	if ForbiddenWordCacheRedis != nil {
		go ForbiddenWordCacheRedis.WarmUp()
	}
	
	// 预热热门话题
	if TopicCacheRedis != nil {
		go TopicCacheRedis.WarmUpHotTopics(100)
	}
	
	slog.Info("Cache warm up initiated")
}

// ClearAll 清除所有缓存
func (m *CacheManager) ClearAll() {
	if !m.redisEnabled {
		return
	}
	
	slog.Info("Clearing all cache...")
	
	client := redis.GetRedisClient()
	if client == nil {
		return
	}
	
	// 使用SCAN命令清除所有带前缀的key
	ctx := context.Background()
	iter := client.Scan(ctx, 0, redis.KeyPrefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		client.Del(ctx, key)
	}
	
	if err := iter.Err(); err != nil {
		slog.Warn("Clear cache failed", "error", err)
	} else {
		slog.Info("All cache cleared")
	}
}