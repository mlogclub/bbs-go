package cache

import (
	"sync"
	"time"

	"github.com/goburrow/cache"
)

// OIDCLoginStateData is server-side only. In particular, the PKCE verifier and
// nonce must never be sent to the browser callback.
type OIDCLoginStateData struct {
	ProviderKey string
	Redirect    string
	Nonce       string
	Verifier    string
}

type oidcLoginStateCache struct {
	cache cache.Cache
	mu    sync.Mutex
}

var OIDCLoginStateCache = &oidcLoginStateCache{cache: cache.New(cache.WithMaximumSize(10000), cache.WithExpireAfterWrite(10*time.Minute))}

func (c *oidcLoginStateCache) Put(state string, data *OIDCLoginStateData) { c.cache.Put(state, data) }

// Take consumes the state before exchanging the code, preventing replay.
func (c *oidcLoginStateCache) Take(state string) *OIDCLoginStateData {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, found := c.cache.GetIfPresent(state)
	if !found {
		return nil
	}
	c.cache.Invalidate(state)
	return value.(*OIDCLoginStateData)
}
