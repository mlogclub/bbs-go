package cache

import (
	"github.com/goburrow/cache"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	indexArticleCacheKey = "index_article_cache"
)

type articleCache struct {
	indexArticlesCache cache.LoadingCache
}

func newArticleCache() *articleCache {
	return &articleCache{
		indexArticlesCache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articles, err := repositories.ArticleRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.ArticleStatusPublished).
					Order("id desc").Size(20))
				if err != nil {
					logrus.Error(err)
				} else {
					value = articles
				}
				return
			},
			cache.WithMaximumSize(10),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (this *articleCache) GetIndexList() []model.Article {
	val, err := this.indexArticlesCache.Get(indexArticleCacheKey)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]model.Article)
	}
	return nil
}

func (this *articleCache) InvalidateIndexList() {
	this.indexArticlesCache.Invalidate(indexArticleCacheKey)
}
