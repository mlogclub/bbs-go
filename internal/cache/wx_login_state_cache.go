package cache

import (
	"time"

	"github.com/goburrow/cache"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
)

type wxLoginStateCache struct {
	cache cache.Cache
}

type WxLoginStateData struct {
	Redirect string
	Bind     bool // 表明当前是不是绑定流程
	UserInfo oauth.UserInfo
}

var WxLoginStateCache = newWxLoginStateCache()

func newWxLoginStateCache() *wxLoginStateCache {
	return &wxLoginStateCache{
		cache: cache.New(
			cache.WithMaximumSize(10000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *wxLoginStateCache) Get(state string) *WxLoginStateData {
	val, found := c.cache.GetIfPresent(state)
	if !found {
		return nil
	}
	return val.(*WxLoginStateData)
}

func (c *wxLoginStateCache) Put(state string, data *WxLoginStateData) {
	c.cache.Put(state, data)
}
