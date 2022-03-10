package gocache

import (
	"fmt"
	"gocache/lru"
	"sync"
)

// 加载 缓存的函数

type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义一个回调函数，获取key不到时候使用

type GetterFunc func(key string) ([]byte, error)

// 回调函数实现 加载缓存的函数

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 操作主体，负责与用户交互

type Group struct {
	name      string // 缓存名称 缓存的命名空间
	getter    Getter // 回调函数
	mainCache cache  // 缓存操作
}

var (
	rw     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, size int64, getter Getter) *Group {
	if getter == nil {
		panic("getter failed")
	}
	rw.Lock()
	defer rw.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: size, lru:lru.New(size, nil)},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	rw.RLock()
	g := groups[name]
	rw.RUnlock()
	return g
}

func (g *Group) Get(key string) (value ByteView, ok error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if get, ok := g.mainCache.get(key); ok {
		return get, nil
	}
	return g.load(key)

}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocal(key)
}

// 使用回调函数获取到缓存
func (g *Group) getLocal(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	view := ByteView{b: cloneBytes(bytes)}
	// 将获取到缓存 写入缓存系统
	g.populateCache(key, view)
	return view, nil
}

// 将获取到缓存 写入缓存系统
func (g *Group) populateCache(key string, view ByteView) {
	g.mainCache.add(key, view)
}
