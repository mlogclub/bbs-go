package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/redis"
	"bbs-go/internal/pkg/simple/common/dates"
	"bbs-go/internal/pkg/simple/sqls"
	"bbs-go/internal/repositories"
)

type redisUserCache struct {
	cache redis.Cache
}

var UserCacheRedis = newRedisUserCache()

func newRedisUserCache() *redisUserCache {
	cache := redis.NewRedisCache()
	if cache == nil {
		return nil
	}
	return &redisUserCache{
		cache: cache,
	}
}

func (c *redisUserCache) Get(userId int64) *models.User {
	if userId <= 0 {
		return nil
	}
	
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getFromDB(userId)
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.UserCacheKey, userId)
	
	// 使用防缓存击穿获取数据
	result, err := c.cache.GetWithBreakdown(ctx, key, redis.UserCacheExpire, func() (interface{}, error) {
		return c.getFromDB(userId), nil
	})
	
	if err != nil {
		slog.Warn("Redis get user failed, fallback to database", "userId", userId, "error", err)
		return c.getFromDB(userId)
	}
	
	if result == nil {
		return nil
	}
	
	user, ok := result.(*models.User)
	if !ok {
		// 类型断言失败，从数据库重新加载
		return c.getFromDB(userId)
	}
	
	return user
}

func (c *redisUserCache) getFromDB(userId int64) *models.User {
	if c == nil {
		return repositories.UserRepository.Get(sqls.DB(), userId)
	}
	return repositories.UserRepository.Get(sqls.DB(), userId)
}

func (c *redisUserCache) Invalidate(userId int64) {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	key := fmt.Sprintf(redis.UserCacheKey, userId)
	
	err := c.cache.Del(ctx, key)
	if err != nil {
		slog.Warn("Redis delete user cache failed", "userId", userId, "error", err)
	}
}

func (c *redisUserCache) GetScoreRank() []models.User {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getScoreRankFromDB()
	}
	
	ctx := context.Background()
	
	var users []models.User
	err := c.cache.GetJSON(ctx, redis.UserScoreRankKey, &users)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		users = c.getScoreRankFromDB()
		if len(users) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), redis.UserScoreRankKey, users, redis.RankCacheExpire)
			}()
		}
	}
	
	return users
}

func (c *redisUserCache) getScoreRankFromDB() []models.User {
	return repositories.UserRepository.Find(sqls.DB(), sqls.NewCnd().Desc("score").Limit(10))
}

func (c *redisUserCache) GetCheckInRank() []models.CheckIn {
	if c == nil || c.cache == nil || !redis.IsRedisEnabled() {
		return c.getCheckInRankFromDB()
	}
	
	ctx := context.Background()
	today := dates.GetDay(time.Now())
	key := fmt.Sprintf(redis.UserCheckInRankKey, today)
	
	var checkIns []models.CheckIn
	err := c.cache.GetJSON(ctx, key, &checkIns)
	if err != nil {
		// 缓存未命中或出错，从数据库加载
		checkIns = c.getCheckInRankFromDB()
		if len(checkIns) > 0 {
			// 异步写入缓存
			go func() {
				c.cache.SetJSON(context.Background(), key, checkIns, time.Hour)
			}()
		}
	}
	
	return checkIns
}

func (c *redisUserCache) getCheckInRankFromDB() []models.CheckIn {
	today := dates.GetDay(time.Now())
	return repositories.CheckInRepository.Find(sqls.DB(),
		sqls.NewCnd().Eq("latest_day_name", today).Asc("update_time").Limit(10))
}

func (c *redisUserCache) RefreshCheckInRank() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	today := dates.GetDay(time.Now())
	key := fmt.Sprintf(redis.UserCheckInRankKey, today)
	
	checkIns := c.getCheckInRankFromDB()
	err := c.cache.SetJSON(ctx, key, checkIns, time.Hour)
	if err != nil {
		slog.Warn("Redis refresh check-in rank failed", "error", err)
	}
}

func (c *redisUserCache) RefreshScoreRank() {
	if !redis.IsRedisEnabled() {
		return
	}
	
	ctx := context.Background()
	users := c.getScoreRankFromDB()
	err := c.cache.SetJSON(ctx, redis.UserScoreRankKey, users, redis.RankCacheExpire)
	if err != nil {
		slog.Warn("Redis refresh score rank failed", "error", err)
	}
}

// 批量预热用户缓存
func (c *redisUserCache) WarmUp(userIds []int64) {
	if !redis.IsRedisEnabled() || len(userIds) == 0 {
		return
	}
	
	ctx := context.Background()
	pipe := c.cache.Pipeline()
	
	for _, userId := range userIds {
		user := c.getFromDB(userId)
		if user != nil {
			key := fmt.Sprintf(redis.UserCacheKey, userId)
			// 使用JSON序列化设置缓存
			if data, err := json.Marshal(user); err == nil {
				pipe.Set(ctx, key, data, redis.UserCacheExpire)
			}
		}
	}
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.Warn("Redis warm up user cache failed", "error", err)
	}
}