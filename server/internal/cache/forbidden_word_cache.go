package cache

import (
	"log/slog"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"

	"github.com/goburrow/cache"

	"bbs-go/internal/pkg/simple/sqls"
)

type forbiddenWordCache struct {
	cache cache.LoadingCache
}

var ForbiddenWordCache = newForbiddenWordCache()

func newForbiddenWordCache() *forbiddenWordCache {
	return &forbiddenWordCache{
		cache: cache.NewLoadingCache(
			func(_ cache.Key) (value cache.Value, e error) {
				value = repositories.ForbiddenWordRepository.Find(sqls.DB(), sqls.NewCnd())
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *forbiddenWordCache) Get() []models.ForbiddenWord {
	val, err := c.cache.Get("_")
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	}
	return val.([]models.ForbiddenWord)
}

func (c *forbiddenWordCache) Invalidate() {
	c.cache.Invalidate("_")
}
