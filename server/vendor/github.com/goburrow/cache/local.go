package cache

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Default maximum number of cache entries.
	defaultCapacity = 1<<15 - 1
	maximumCapacity = 1<<31 - 1
	// Buffer size of entry channels
	chanBufSize = 1
	// Maximum number of entries to be drained in a single clean up.
	drainMax       = 16
	drainThreshold = 64
)

// currentTime is an alias for time.Now, used for testing.
var currentTime = time.Now

// localCache is an asynchronous LRU cache.
type localCache struct {
	// user configurations
	policyName        string
	expireAfterAccess time.Duration
	refreshAfterWrite time.Duration

	onInsertion Func
	onRemoval   Func

	loader LoaderFunc
	stats  StatsCounter

	// internal data structure
	cap   int
	cache cache

	entries     policy
	addEntry    chan *entry
	hitEntry    chan *list.Element
	deleteEntry chan *list.Element

	// readCount is a counter of the number of reads since the last write.
	readCount int32

	// for closing routines created by this cache.
	closeMu sync.Mutex
	closeCh chan struct{}
}

// newLocalCache returns a default localCache.
// init must be called before this cache can be used.
func newLocalCache() *localCache {
	return &localCache{
		cap: defaultCapacity,
		cache: cache{
			data: make(map[Key]*list.Element),
		},
		stats: &statsCounter{},
	}
}

// init initializes cache replacement policy after all user configuration properties are set.
func (c *localCache) init() {
	c.entries = newPolicy(c.policyName)
	c.entries.init(&c.cache, c.cap)

	c.addEntry = make(chan *entry, chanBufSize)
	c.hitEntry = make(chan *list.Element, chanBufSize)
	c.deleteEntry = make(chan *list.Element, chanBufSize)

	c.closeCh = make(chan struct{})
	go c.processEntries()
}

// Close implements io.Closer and always returns a nil error.
func (c *localCache) Close() error {
	c.closeMu.Lock()
	defer c.closeMu.Unlock()
	if c.closeCh != nil {
		c.closeCh <- struct{}{}
		// Wait for the goroutine to close this channel
		// (should use sync.WaitGroup or a new channel instead?)
		<-c.closeCh
		c.closeCh = nil
	}
	return nil
}

// GetIfPresent gets cached value from entries list and updates
// last access time for the entry if it is found.
func (c *localCache) GetIfPresent(k Key) (Value, bool) {
	c.cache.mu.RLock()
	el, hit := c.cache.data[k]
	c.cache.mu.RUnlock()
	if !hit {
		c.stats.RecordMisses(1)
		return nil, false
	}
	v := getEntry(el).value
	c.hitEntry <- el
	c.stats.RecordHits(1)
	return v, true
}

// Put adds new entry to entries list.
func (c *localCache) Put(k Key, v Value) {
	c.cache.mu.RLock()
	el, hit := c.cache.data[k]
	c.cache.mu.RUnlock()
	if hit {
		// Update list element value
		getEntry(el).value = v
		c.hitEntry <- el
	} else {
		en := &entry{
			key:   k,
			value: v,
			hash:  sum(k),
		}
		c.addEntry <- en
	}
}

// Invalidate removes the entry associated with key k.
func (c *localCache) Invalidate(k Key) {
	c.cache.mu.RLock()
	el, hit := c.cache.data[k]
	c.cache.mu.RUnlock()
	if hit {
		c.deleteEntry <- el
	}
}

// InvalidateAll resets entries list.
func (c *localCache) InvalidateAll() {
	c.deleteEntry <- nil
}

// Get returns value associated with k or call underlying loader to retrieve value
// if it is not in the cache. The returned value is only cached when loader returns
// nil error.
func (c *localCache) Get(k Key) (Value, error) {
	c.cache.mu.RLock()
	el, hit := c.cache.data[k]
	c.cache.mu.RUnlock()
	if !hit {
		c.stats.RecordMisses(1)
		return c.load(k)
	}
	c.stats.RecordHits(1)
	en := getEntry(el)
	// Check if this entry needs to be refreshed
	if c.refreshAfterWrite > 0 && en.updated.Before(currentTime().Add(-c.refreshAfterWrite)) {
		return c.refresh(en), nil
	}
	v := en.value
	c.hitEntry <- el
	return v, nil
}

// Stats copies cache stats to t.
func (c *localCache) Stats(t *Stats) {
	c.stats.Snapshot(t)
}

func (c *localCache) processEntries() {
	defer close(c.closeCh)
	for {
		select {
		case <-c.closeCh:
			c.removeAll()
			return
		case en := <-c.addEntry:
			c.add(en)
			c.postWriteCleanup()
		case el := <-c.hitEntry:
			c.hit(el)
			c.postReadCleanup()
		case el := <-c.deleteEntry:
			if el == nil {
				c.removeAll()
			} else {
				c.remove(el)
			}
			c.postReadCleanup()
		}
	}
}

func (c *localCache) add(en *entry) {
	en.accessed = currentTime()
	en.updated = en.accessed

	remEn := c.entries.add(en)
	if c.onInsertion != nil {
		c.onInsertion(en.key, en.value)
	}
	if remEn != nil {
		// An entry has been evicted
		c.stats.RecordEviction()
		if c.onRemoval != nil {
			c.onRemoval(remEn.key, remEn.value)
		}
	}
}

// removeAll remove all entries in the cache.
func (c *localCache) removeAll() {
	c.cache.mu.Lock()
	oldData := c.cache.data
	c.cache.data = make(map[Key]*list.Element)
	c.entries.init(&c.cache, c.cap)
	c.cache.mu.Unlock()

	if c.onRemoval != nil {
		for _, el := range oldData {
			en := getEntry(el)
			c.onRemoval(en.key, en.value)
		}
	}
}

// remove removes the given element from the cache and entries list.
// It also calls onRemoval callback if it is set.
func (c *localCache) remove(el *list.Element) {
	en := c.entries.remove(el)

	if en != nil && c.onRemoval != nil {
		c.onRemoval(en.key, en.value)
	}
}

// hit moves the given element to the top of the entries list.
func (c *localCache) hit(el *list.Element) {
	getEntry(el).accessed = currentTime()
	c.entries.hit(el)
}

// load uses current loader to retrieve value for k and adds new
// entry to the cache only if loader returns a nil error.
func (c *localCache) load(k Key) (Value, error) {
	if c.loader == nil {
		panic("loader must be set")
	}
	start := currentTime()
	v, err := c.loader(k)
	loadTime := currentTime().Sub(start)
	if err != nil {
		c.stats.RecordLoadError(loadTime)
		return nil, err
	}
	en := &entry{
		key:   k,
		value: v,
		hash:  sum(k),
	}
	c.addEntry <- en
	c.stats.RecordLoadSuccess(loadTime)
	return v, nil
}

// refresh reloads value for the given key. If loader returns an error,
// that error will be omitted and current value will be returned.
// Otherwise, the function will returns new value and updates the current
// cache entry.
func (c *localCache) refresh(en *entry) Value {
	if c.loader == nil {
		panic("loader must be set")
	}
	start := currentTime()
	newV, err := c.loader(en.key)
	loadTime := currentTime().Sub(start)
	if err != nil {
		c.stats.RecordLoadError(loadTime)
		return en.value
	}
	en.value = newV
	c.addEntry <- en
	c.stats.RecordLoadSuccess(loadTime)
	return newV
}

// postReadCleanup is run after entry access/delete event.
func (c *localCache) postReadCleanup() {
	if atomic.AddInt32(&c.readCount, 1) > drainThreshold {
		atomic.StoreInt32(&c.readCount, 0)
		c.expireEntries()
	}
}

// postWriteCleanup is run after entry add event.
func (c *localCache) postWriteCleanup() {
	atomic.StoreInt32(&c.readCount, 0)
	c.expireEntries()
}

// expireEntries removes expired entries.
func (c *localCache) expireEntries() {
	if c.expireAfterAccess <= 0 {
		return
	}
	expire := currentTime().Add(-c.expireAfterAccess)
	remain := drainMax
	c.entries.walk(func(ls *list.List) {
		for ; remain > 0; remain-- {
			el := ls.Back()
			if el == nil {
				// List is empty
				break
			}
			en := getEntry(el)
			if !en.accessed.Before(expire) {
				// Can break since the entries list is sorted by access time
				break
			}
			c.remove(el)
			c.stats.RecordEviction()
		}
	})
}

// New returns a local in-memory Cache.
func New(options ...Option) Cache {
	c := newLocalCache()
	for _, opt := range options {
		opt(c)
	}
	c.init()
	return c
}

// NewLoadingCache returns a new LoadingCache with given loader function
// and cache options.
func NewLoadingCache(loader LoaderFunc, options ...Option) LoadingCache {
	c := newLocalCache()
	c.loader = loader
	for _, opt := range options {
		opt(c)
	}
	c.init()
	return c
}

// Option add options for default Cache.
type Option func(c *localCache)

// WithMaximumSize returns an Option which sets maximum size for the cache.
// Any non-positive numbers is considered as unlimited.
func WithMaximumSize(size int) Option {
	if size < 0 {
		size = 0
	}
	if size > maximumCapacity {
		size = maximumCapacity
	}
	return func(c *localCache) {
		c.cap = size
	}
}

// WithRemovalListener returns an Option to set cache to call onRemoval for each
// entry evicted from the cache.
func WithRemovalListener(onRemoval Func) Option {
	return func(c *localCache) {
		c.onRemoval = onRemoval
	}
}

// WithExpireAfterAccess returns an option to expire a cache entry after the
// given duration without being accessed.
func WithExpireAfterAccess(d time.Duration) Option {
	return func(c *localCache) {
		c.expireAfterAccess = d
	}
}

// WithRefreshAfterWrite returns an option to refresh a cache entry after the
// given duration. This option is only applicable for LoadingCache.
func WithRefreshAfterWrite(d time.Duration) Option {
	return func(c *localCache) {
		c.refreshAfterWrite = d
	}
}

// WithStatsCounter returns an option which overrides default cache stats counter.
func WithStatsCounter(st StatsCounter) Option {
	return func(c *localCache) {
		c.stats = st
	}
}

// WithPolicy returns an option which set cache policy associated to the given name.
// Supported policies are: lru, slru, tinylfu.
func WithPolicy(name string) Option {
	return func(c *localCache) {
		c.policyName = name
	}
}

// withInsertionListener is used for testing.
func withInsertionListener(onInsertion Func) Option {
	return func(c *localCache) {
		c.onInsertion = onInsertion
	}
}
