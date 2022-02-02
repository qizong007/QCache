package qcache

import (
	lru2 "QCache/qcache/lru"
	"sync"
)

// 对 lru 的进一步封装，主要是做并发控制
type cache struct {
	mu         sync.Mutex
	lru        *lru2.Cache
	cacheBytes int64
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru2.New(c.cacheBytes, nil) // Lazy Initialization
	}
	c.lru.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
