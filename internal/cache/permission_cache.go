package cache

import "sync"

type permissionCache struct {
	mu    sync.RWMutex
	codes map[int64][]string
}

var PermissionCache = newPermissionCache()

func newPermissionCache() *permissionCache {
	return &permissionCache{
		codes: make(map[int64][]string),
	}
}

func (c *permissionCache) Get(userId int64) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	codes, ok := c.codes[userId]
	if !ok {
		return nil, false
	}
	return append([]string(nil), codes...), true
}

func (c *permissionCache) Put(userId int64, codes []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.codes[userId] = append([]string(nil), codes...)
}

func (c *permissionCache) Invalidate(userId int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.codes, userId)
}

func (c *permissionCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.codes = make(map[int64][]string)
}
