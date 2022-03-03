package main

import (
	"fmt"
	"go-cache/lru"
	"log"
)

// OnEvicted 自定义回调接口
func OnEvicted(key string, value lru.Value)  {
	log.Printf("cache evicted, key:%s, value:%v", key, value)
}

type String string
// Len 定义缓存的数据结构，实现Value 接口
func (s String) Len() int{
	return len(s)
}

func main() {


	cache := lru.New(int64(20), OnEvicted)
	cache.Add("test", String("val"))
	cache.Add("zmn", String("val"))
	// 在增加数据时候，内存不够会触发 LRU
	cache.Add("lisi", String("val"))
	cache.Add("python", String("val"))
	fmt.Println(cache.Len())

	if _, ok := cache.Get("zmn"); !ok{
		fmt.Println("not hit cache, key:zmn miss" )
	}

}
