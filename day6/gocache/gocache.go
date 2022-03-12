package gocache

import (
	"fmt"
	"gocache/lru"
	"log"
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
	name       string     // 缓存名称 缓存的命名空间
	getter     Getter     // 回调函数
	mainCache  cache      // 缓存操作
	peerPicker PeerPicker // 新增获取节点的能力
	loader     *FGroup    // 增加防止缓存击穿能力
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
		mainCache: cache{cacheBytes: size, lru: lru.New(size, nil)},
		loader:    &FGroup{m: make(map[string]*call)},
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

// 根据key 获取缓存的方法
func (g *Group) load(key string) (value ByteView, err error) {

	// 确保统一时间 调用远程服务，只有一次
	view, e := g.loader.DO(key, func() (interface{}, error) {
		if g.peerPicker != nil {
			// 查看key 在那个节点，区分调用方式
			if peer, ok := g.peerPicker.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					log.Println("[GeeCache] load cache from remote peer: ", peer)
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		log.Println("[GeeCache]load cache from local peer ")
		return g.getLocal(key)
	})

	if e == nil {
		return view.(ByteView), nil
	}
	return
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

// RegisterPeerPicker 给group 注入 picker
func (g *Group) RegisterPeerPicker(picker PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker called more than once ")
	}
	g.peerPicker = picker
}

// 从远程节点获取数据的方法
func (g *Group) getFromPeer(getter PeerGetter, key string) (ByteView, error) {
	get, err := getter.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: get}, nil
}
