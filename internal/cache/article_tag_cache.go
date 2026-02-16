package cache

import (
	"errors"
	"time"

	"github.com/goburrow/cache"
	"github.com/mlogclub/simple/sqls"

	"bbs-go/internal/repositories"
)

type articleTagCache struct {
	cache cache.LoadingCache
}

var ArticleTagCache = newArticleTagCache()

func newArticleTagCache() *articleTagCache {
	return &articleTagCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				articleTags := repositories.ArticleTagRepository.FindByArticleId(sqls.DB(), key2Int64(key))
				if len(articleTags) > 0 {
					var tagIds []int64
					for _, articleTag := range articleTags {
						tagIds = append(tagIds, articleTag.TagId)
					}
					value = tagIds
				} else {
					e = errors.New("文章没标签")
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *articleTagCache) Get(articleId int64) []int64 {
	val, err := c.cache.Get(articleId)
	if err != nil {
		return nil
	}
	if val != nil {
		return val.([]int64)
	}
	return nil
}

func (c *articleTagCache) Invalidate(articleId int64) {
	c.cache.Invalidate(articleId)
}
