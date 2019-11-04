package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/repositories"
)

type articleTagCache struct {
	cache cache.LoadingCache
}

var ArticleTagCache = newArticleTagCache()

func newArticleTagCache() *articleTagCache {
	return &articleTagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articleTags := repositories.ArticleTagRepository.FindByArticleId(simple.DB(), Key2Int64(key))
				if len(articleTags) > 0 {
					var tagIds []int64
					for _, articleTag := range articleTags {
						tagIds = append(tagIds, articleTag.TagId)
					}
					value = tagIds
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (this *articleTagCache) Get(articleId int64) []int64 {
	val, err := this.cache.Get(articleId)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]int64)
	}
	return nil
}

func (this *articleTagCache) Invalidate(articleId int64) {
	this.cache.Invalidate(articleId)
}
