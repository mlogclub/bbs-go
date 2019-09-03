package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

var (
	recommendCacheKey = "recommend_articles_cache"
)

var ArticleCache = newArticleCache()

type articleCache struct {
	recommendCache     cache.LoadingCache
}

func newArticleCache() *articleCache {
	return &articleCache{
		recommendCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articles, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(),
					simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).Order("id desc").Size(50))
				if err != nil {
					logrus.Error(err)
				} else {
					value = articles
				}
				return
			},
			cache.WithMaximumSize(1),
			cache.WithExpireAfterAccess(10*time.Minute),
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
