package main

import (
	"fmt"
	"gocache"
	"log"
	"reflect"
)

func TestGetter() {
	getterFunc := gocache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	v, _ := getterFunc.Get("key")

	if !reflect.DeepEqual(v, []byte("key")) {
		log.Fatal("callback faild")
	} else {
		log.Println("success")
	}
}

// 测试组

func testGroup() {

	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	loadCounts := make(map[string]int, 10)

	group := gocache.NewGroup("score", 15, gocache.GetterFunc(func(key string) ([]byte, error) {

		log.Printf("search %s from db", key)
		if s, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key]++
			return []byte(s), nil
		}
		return nil, fmt.Errorf("key not found from db")
	}))

	for k, _ := range db {
		_, err := group.Get(k)
		if err != nil {
			log.Printf("key get failed: %s", err)
		}
		if loadCounts[k] > 1 {
			log.Printf("key %s cache miss, from db >1", k)
		}
	}

	for k, _ := range db {
		_, err := group.Get(k)
		if err != nil {
			log.Printf("key get failed: %s", err)
		}
		if loadCounts[k] > 1 {
			log.Printf("key %s cache miss, from db >1", k)
		}
	}

	get, err := group.Get("test")
	if err != nil{
		log.Fatalf("test err: %s", err)
	}else {
		log.Printf("get test:%s", get)
	}

}

func main() {
	//TestGetter()

	testGroup()
}
