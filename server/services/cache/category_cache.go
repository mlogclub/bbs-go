package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

type categoryCache struct {
	cache cache.LoadingCache
}

var CategoryCache = newCategoryCache()

func newCategoryCache() *categoryCache {
	return &categoryCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = repositories.CategoryRepository.Get(simple.DB(), Key2Int64(key))
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (this *categoryCache) Get(categoryId int64) *model.Category {
	val, err := this.cache.Get(categoryId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if val != nil {
		return val.(*model.Category)
	}
	return nil
}

func (this *categoryCache) Invalidate(categoryId int64) {
	this.cache.Invalidate(categoryId)
}
