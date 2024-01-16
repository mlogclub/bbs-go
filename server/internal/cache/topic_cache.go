package cache

import (
	"bbs-go/internal/models/constants"
	"errors"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var (
	topicRecommendCacheKey = "recommend_topics_cache"
)

var TopicCache = newTopicCache()

type topicCache struct {
	recommendCache cache.LoadingCache
}

func newTopicCache() *topicCache {
	return &topicCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.TopicRepository.Find(sqls.DB(),
					sqls.NewCnd().Eq("status", constants.StatusOk).Desc("id").Limit(50))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (c *topicCache) GetRecommendTopics() []models.Topic {
	val, err := c.recommendCache.Get(topicRecommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]models.Topic)
	}
	return nil
}

func (c *topicCache) InvalidateRecommend() {
	c.recommendCache.Invalidate(topicRecommendCacheKey)
}
