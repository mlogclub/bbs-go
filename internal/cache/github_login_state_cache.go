package cache

import (
	"time"

	"bbs-go/internal/pkg/github"

	"github.com/goburrow/cache"
)

type githubLoginStateCache struct {
	cache cache.Cache
}

type GithubLoginStateData struct {
	Redirect string
	Bind     bool // 表明当前是不是绑定流程
	UserInfo *github.GithubUserInfo
}

var GithubLoginStateCache = newGithubLoginStateCache()

func newGithubLoginStateCache() *githubLoginStateCache {
	return &githubLoginStateCache{
		cache: cache.New(
			cache.WithMaximumSize(10000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *githubLoginStateCache) Get(state string) *GithubLoginStateData {
	val, found := c.cache.GetIfPresent(state)
	if !found {
		return nil
	}
	return val.(*GithubLoginStateData)
}

func (c *githubLoginStateCache) Put(state string, data *GithubLoginStateData) {
	c.cache.Put(state, data)
}
