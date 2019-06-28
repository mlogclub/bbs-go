package cache

import (
	"github.com/goburrow/cache"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"time"
)

type articleTagCache struct {
	cache cache.LoadingCache
}

func newArticleTagCache() *articleTagCache {
	return &articleTagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articleTags, err := repositories.NewArticleTagRepository().GetByArticleId(simple.GetDB(), Key2Int64(key))
				if err != nil {
					logrus.Error(err)
					return
				}
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
