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

type redisArticleTagCache struct {
	cache redis.Cache
}

var ArticleTagCacheRedis = newRedisArticleTagCache()

func newRedisArticleTagCache() *redisArticleTagCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisArticleTagCache{
		cache: cache,
	}
}

func (c *redisArticleTagCache) Get(articleId int64) []int64 {
	if articleId <= 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getTagIdsFromDB(articleId)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.ArticleTagCacheKey, articleId)
	
	// 先尝试获取标签ID数组（为了兼容性）
	var tagIds []int64
	err := c.cache.GetJSON(ctx, key, &tagIds)
	if err == nil {
		return tagIds
	}
	
	// 如果失败，从数据库加载
	articleTags := c.getFromDB(articleId)
	if len(articleTags) > 0 {
		// 提取标签ID
		for _, articleTag := range articleTags {
			tagIds = append(tagIds, articleTag.TagId)
		}
		
		// 异步写入缓存（存储标签ID数组）
		go func() {
			c.cache.SetJSON(context.Background(), key, tagIds, redis.ArticleTagCacheExpire)
		}()
		
		return tagIds
	}
	
	// 设置空值缓存防穿透
	go func() {
		redisCache, ok := c.cache.(*redis.RedisCache)
		if ok {
			redisCache.SetNullValue(context.Background(), key, redis.ArticleTagCacheExpire)
		}
	}()
	
	return nil
}

func (c *redisArticleTagCache) GetArticleTags(articleId int64) []models.ArticleTag {
	if articleId <= 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(articleId)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.ArticleTagCacheKey+"_full", articleId)
	
	var articleTags []models.ArticleTag
	err := c.cache.GetJSON(ctx, key, &articleTags)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		articleTags = c.getFromDB(articleId)
		if len(articleTags) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), key, articleTags, redis.ArticleTagCacheExpire)
			}()
		}
	}
	
	return articleTags
}

func (c *redisArticleTagCache) getTagIdsFromDB(articleId int64) []int64 {
	if c == nil {
		articleTags := repositories.ArticleTagRepository.Find(sqls.DB(), 
			sqls.NewCnd().Eq("article_id", articleId))
		var tagIds []int64
		for _, articleTag := range articleTags {
			tagIds = append(tagIds, articleTag.TagId)
		}
		return tagIds
	}
	articleTags := c.getFromDB(articleId)
	var tagIds []int64
	for _, articleTag := range articleTags {
		tagIds = append(tagIds, articleTag.TagId)
	}
	return tagIds
}

func (c *redisArticleTagCache) getFromDB(articleId int64) []models.ArticleTag {
	if c == nil {
		return repositories.ArticleTagRepository.Find(sqls.DB(), 
			sqls.NewCnd().Eq("article_id", articleId))
	}
	return repositories.ArticleTagRepository.Find(sqls.DB(), 
		sqls.NewCnd().Eq("article_id", articleId))
}

func (c *redisArticleTagCache) Invalidate(articleId int64) {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.ArticleTagCacheKey, articleId)
	
	err := c.cache.Del(ctx, key)
	if err != nil {
		slog.Warn("Redis delete article tag cache failed", "articleId", articleId, "error", err)
	}
}

func (c *redisArticleTagCache) Set(articleId int64, articleTags []models.ArticleTag) {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.ArticleTagCacheKey, articleId)
	
	if len(articleTags) > 0 {
		err := c.cache.SetJSON(ctx, key, articleTags, redis.ArticleTagCacheExpire)
		if err != nil {
			slog.Warn("Redis set article tag cache failed", "articleId", articleId, "error", err)
		}
	} else {
		// 设置空值缓存
		redisCache, ok := c.cache.(*redis.RedisCache)
		if ok {
			err := redisCache.SetNullValue(ctx, key, redis.ArticleTagCacheExpire)
			if err != nil {
				slog.Warn("Redis set null article tag cache failed", "articleId", articleId, "error", err)
			}
		}
	}
}