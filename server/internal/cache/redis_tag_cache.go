package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisTagCache struct {
	cache redis.Cache
}

var TagCacheRedis = newRedisTagCache()

func newRedisTagCache() *redisTagCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisTagCache{
		cache: cache,
	}
}

func (c *redisTagCache) Get(tagId int64) *models.Tag {
	if tagId <= 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(tagId)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.TagCacheKey, tagId)
	
	// 使用防缓存击穿获取数据
	result, err := c.cache.GetWithBreakdown(ctx, key, redis.TagCacheExpire, func() (interface{}, error) {
		return c.getFromDB(tagId), nil
	})
	
	if err != nil {
		slog.Warn("Redis get tag failed, fallback to database", "tagId", tagId, "error", err)
		return c.getFromDB(tagId)
	}
	
	if result == nil {
		return nil
	}
	
	tag, ok := result.(*models.Tag)
	if !ok {
		// 类型断言失败，从数据库重新加载
		return c.getFromDB(tagId)
	}
	
	return tag
}

func (c *redisTagCache) getFromDB(tagId int64) *models.Tag {
	if c == nil {
		return repositories.TagRepository.Get(sqls.DB(), tagId)
	}
	return repositories.TagRepository.Get(sqls.DB(), tagId)
}

func (c *redisTagCache) GetAll() []models.Tag {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getAllFromDB()
	}
	
	ctx := context.Background()
	
	var tags []models.Tag
	err := c.cache.GetJSON(ctx, redis.TagAllCacheKey, &tags)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		tags = c.getAllFromDB()
		if len(tags) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), redis.TagAllCacheKey, tags, redis.TagCacheExpire)
			}()
		}
	}
	
	return tags
}

func (c *redisTagCache) getAllFromDB() []models.Tag {
	if c == nil {
		return repositories.TagRepository.Find(sqls.DB(), sqls.NewCnd())
	}
	return repositories.TagRepository.Find(sqls.DB(), sqls.NewCnd())
}

func (c *redisTagCache) GetList(tagIds []int64) (tags []models.Tag) {
	if c == nil || len(tagIds) == 0 {
		return nil
	}
	for _, tagId := range tagIds {
		tag := c.Get(tagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (c *redisTagCache) Invalidate(tagId int64) {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	keys := []string{
		fmt.Sprintf(redis.TagCacheKey, tagId),
		redis.TagAllCacheKey, // 也清除全量缓存
	}
	
	err := c.cache.Del(ctx, keys...)
	if err != nil {
		slog.Warn("Redis delete tag cache failed", "tagId", tagId, "error", err)
	}
}

func (c *redisTagCache) InvalidateAll() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	err := c.cache.Del(ctx, redis.TagAllCacheKey)
	if err != nil {
		slog.Warn("Redis delete all tag cache failed", "error", err)
	}
}

func (c *redisTagCache) Set(tagId int64, tag *models.Tag) {
	if tag == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.TagCacheKey, tagId)
	
	err := c.cache.SetJSON(ctx, key, tag, redis.TagCacheExpire)
	if err != nil {
		slog.Warn("Redis set tag cache failed", "tagId", tagId, "error", err)
	}
	
	// 同时刷新全量缓存
	go c.RefreshAll()
}

func (c *redisTagCache) RefreshAll() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	tags := c.getAllFromDB()
	
	err := c.cache.SetJSON(ctx, redis.TagAllCacheKey, tags, redis.TagCacheExpire)
	if err != nil {
		slog.Warn("Redis refresh all tag cache failed", "error", err)
	}
}

// 批量预热标签缓存
func (c *redisTagCache) WarmUp() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	tags := c.getAllFromDB()
	if len(tags) == 0 {
		return
	}
	
	ctx := context.Background()
	pipe := c.cache.Pipeline()
	
	// 缓存所有标签
	for _, tag := range tags {
		key := fmt.Sprintf(redis.TagCacheKey, tag.Id)
		// 使用JSON序列化设置缓存
		if data, err := json.Marshal(&tag); err == nil {
			pipe.Set(ctx, key, data, redis.TagCacheExpire)
		}
	}
	
	// 缓存全量标签
	if data, err := json.Marshal(tags); err == nil {
		pipe.Set(ctx, redis.TagAllCacheKey, data, redis.TagCacheExpire)
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Warn("Redis warm up tag cache failed", "error", err)
	} else {
		slog.Info("Tag cache warmed up", "count", len(tags))
	}
}