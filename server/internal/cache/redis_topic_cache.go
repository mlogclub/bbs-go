package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisTopicCache struct {
	cache redis.Cache
}

var TopicCacheRedis = newRedisTopicCache()

func newRedisTopicCache() *redisTopicCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisTopicCache{
		cache: cache,
	}
}

func (c *redisTopicCache) Get(topicId int64) *models.Topic {
	if topicId <= 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(topicId)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.TopicCacheKey, topicId)
	
	// 使用防缓存击穿获取数据
	result, err := c.cache.GetWithBreakdown(ctx, key, redis.TopicCacheExpire, func() (interface{}, error) {
		return c.getFromDB(topicId), nil
	})
	
	if err != nil {
		slog.Warn("Redis get topic failed, fallback to database", "topicId", topicId, "error", err)
		return c.getFromDB(topicId)
	}
	
	if result == nil {
		return nil
	}
	
	topic, ok := result.(*models.Topic)
	if !ok {
		// 类型断言失败，从数据库重新加载
		return c.getFromDB(topicId)
	}
	
	return topic
}

func (c *redisTopicCache) getFromDB(topicId int64) *models.Topic {
	if c == nil {
		return repositories.TopicRepository.Get(sqls.DB(), topicId)
	}
	return repositories.TopicRepository.Get(sqls.DB(), topicId)
}

func (c *redisTopicCache) Invalidate(topicId int64) {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.TopicCacheKey, topicId)
	
	err := c.cache.Del(ctx, key)
	if err != nil {
		slog.Warn("Redis delete topic cache failed", "topicId", topicId, "error", err)
	}
}

func (c *redisTopicCache) Set(topicId int64, topic *models.Topic) {
	if c == nil || c.cache == nil || topic == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.TopicCacheKey, topicId)
	
	err := c.cache.SetJSON(ctx, key, topic, redis.TopicCacheExpire)
	if err != nil {
		slog.Warn("Redis set topic cache failed", "topicId", topicId, "error", err)
	}
}

func (c *redisTopicCache) GetRecommendTopics() []models.Topic {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getRecommendTopicsFromDB()
	}
	
	ctx := context.Background()
	recommendKey := "bbs:recommend_topics"
	
	var topics []models.Topic
	err := c.cache.GetJSON(ctx, recommendKey, &topics)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		topics = c.getRecommendTopicsFromDB()
		if len(topics) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), recommendKey, topics, 30*redis.TopicCacheExpire)
			}()
		}
	}
	
	return topics
}

func (c *redisTopicCache) getRecommendTopicsFromDB() []models.Topic {
	return repositories.TopicRepository.Find(sqls.DB(),
		sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(50))
}

func (c *redisTopicCache) InvalidateRecommend() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	recommendKey := "bbs:recommend_topics"
	err := c.cache.Del(ctx, recommendKey)
	if err != nil {
		slog.Warn("Redis delete recommend topics cache failed", "error", err)
	}
}

// 批量预热话题缓存
func (c *redisTopicCache) WarmUpHotTopics(limit int) {
	if !redis.IsRedisEnabled() {
		return
	}
	
	// 获取热门话题（按浏览量排序）
	topics := repositories.TopicRepository.Find(sqls.DB(), 
		sqls.NewCnd().Desc("view_count").Limit(limit))
	
	if len(topics) == 0 {
		return
	}
	
	ctx := context.Background()
	pipe := c.cache.Pipeline()
	
	for _, topic := range topics {
		key := fmt.Sprintf(redis.TopicCacheKey, topic.Id)
		// 使用JSON序列化设置缓存
		if data, err := json.Marshal(&topic); err == nil {
			pipe.Set(ctx, key, data, redis.TopicCacheExpire)
		}
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Warn("Redis warm up topic cache failed", "error", err)
	} else {
		slog.Info("Topic cache warmed up", "count", len(topics))
	}
}