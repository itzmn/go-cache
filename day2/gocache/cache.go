package gocache

import (
	"gocache/lru"
	"sync"
)

// 缓存模块  并发操作类，包装底层操作

type cache struct {
	mu         sync.Mutex // 互斥锁
	lru        *lru.Cache // 缓存的底层逻辑
	cacheBytes int64      // 缓存大小
}

// 增加数据
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
