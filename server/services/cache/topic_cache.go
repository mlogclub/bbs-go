package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var TopicCache = newTopicCache()

type topicCache struct {
	recommendCache cache.LoadingCache
}

func newTopicCache() *topicCache {
	return &topicCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.TopicRepository.Find(simple.DB(),
					simple.NewSqlCnd().Where("status = ?", model.TopicStatusOk).Desc("id").Limit(50))
				return
			},
			cache.WithMaximumSize(10),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (this *topicCache) GetRecommendTopics() []model.Topic {
	val, err := this.recommendCache.Get(recommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Topic)
	}
	return nil
}

func (this *topicCache) InvalidateRecommend() {
	this.recommendCache.Invalidate(recommendCacheKey)
}
