package cache

import (
	"time"

	"github.com/goburrow/cache"
	"bbs-go/internal/pkg/google"
)

type googleLoginStateCache struct {
	cache cache.Cache
}

type GoogleLoginStateData struct {
	Redirect string
	Bind     bool // 表明当前是不是绑定流程
	UserInfo *google.GoogleUserInfo
}

var GoogleLoginStateCache = newGoogleLoginStateCache()

func newGoogleLoginStateCache() *googleLoginStateCache {
	return &googleLoginStateCache{
		cache: cache.New(
			cache.WithMaximumSize(10000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *googleLoginStateCache) Get(state string) *GoogleLoginStateData {
	val, found := c.cache.GetIfPresent(state)
	if !found {
		return nil
	}
	return val.(*GoogleLoginStateData)
}

func (c *googleLoginStateCache) Put(state string, data *GoogleLoginStateData) {
	c.cache.Put(state, data)
}
