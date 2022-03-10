package lru

import "container/list"

// LRU cache

type Cache struct {
	maxBytes  int64                         // 缓存最大内存
	nBytes    int64                         // 现在使用的内存
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 缓存的内容是双向链表的节点指针
	OnEvicted func(key string, value Value) // 缓存删除时候的回调函数
}

// 双向链表的 数据类型
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 返回数据内存大小
}

func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {

	return &Cache{
		cache:     make(map[string]*list.Element, 0),
		OnEvicted: OnEvicted,
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
	}

}

// Get 根据key 查找缓存
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		// 查找到元素,将元素位置移动到双向链表的正面
		c.ll.MoveToFront(element)
		// 返回存储的数据节点
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// RemoveOldest 清除缓存多余的数据，链表反面数据，反面为访问较少数据
func (c *Cache) RemoveOldest() {
	back := c.ll.Back()
	if back != nil {
		// 从map中删除
		entry := back.Value.(*entry)
		delete(c.cache, entry.key)
		c.ll.Remove(back)
		// 更新当前使用的内存
		c.nBytes -= int64(len(entry.key)) + int64(entry.value.Len())

		// 清除数据的回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}
}

// Add 增加和更新元素方法
func (c *Cache) Add(key string, value Value) bool {

	addSize := int64(len(key)) + int64(value.Len())
	if addSize > c.maxBytes {
		return false
	}
	// 拿到可以的内存
	for (c.nBytes + addSize) > c.maxBytes {
		c.RemoveOldest()
	}

	if element, ok := c.cache[key]; ok {
		// 如果节点已经存在, 将节点移动到正面
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		// 更新使用内存
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		// 更新值
		kv.value = value
	} else {
		// 新增节点
		front := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = front
		c.nBytes += addSize
	}

	return true
}

// Len 获得缓存中的数据行数
func (c *Cache) Len() (int,int64) {

	return c.ll.Len(), c.nBytes
}
