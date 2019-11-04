package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var (
	recommendCacheKey = "recommend_articles_cache"
)

var ArticleCache = newArticleCache()

type articleCache struct {
	recommendCache cache.LoadingCache
}

func newArticleCache() *articleCache {
	return &articleCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.ArticleRepository.Find(simple.DB(),
					simple.NewSqlCnd().Where("status = ?", model.ArticleStatusPublished).Desc("id").Limit(50))
				return
			},
			cache.WithMaximumSize(1),
			cache.WithRefreshAfterWrite(30*time.Minute),
		),
	}
}

func (this *articleCache) GetRecommendArticles() []model.Article {
	val, err := this.recommendCache.Get(recommendCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Article)
	}
	return nil
}

func (this *articleCache) InvalidateRecommend() {
	this.recommendCache.Invalidate(recommendCacheKey)
}
