package cache

import (
	"github.com/goburrow/cache"
	"time"
)

var (
	Cache                cache.Cache
	IndexArticleCacheKey = "index_article_cache"
)

func init() {

	Cache = cache.New(
		cache.WithMaximumSize(1000),
		cache.WithExpireAfterAccess(30*time.Minute))

}
