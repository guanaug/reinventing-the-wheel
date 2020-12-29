package group

import (
	"gerocache/lru"
	"sync"
)

type ByteView struct {
	data []byte
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) String() string {
	return string(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return b.data
}

type cache struct {
	cache *lru.Cache
	mu    sync.Mutex
}

func newCache(maxBytes int64) *cache {
	return &cache{cache: lru.New(maxBytes)}
}

func (c *cache) get(key string) (view ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.cache.Get(key); ok {
		return val.(ByteView), ok
	}

	return
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Set(key, value)
}
