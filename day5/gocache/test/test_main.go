package main

import (
	"fmt"
	"gocache"
	"log"
	"net/http"
	"reflect"
	"sort"
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
	if err != nil {
		log.Fatalf("test err: %s", err)
	} else {
		log.Printf("get test:%s", get)
	}

}

// day3 测试HTTP请求 中使用缓存
func testHttp() {

	var db = map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}

	gocache.NewGroup("score", 2<<10, gocache.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("search %s from db", key)
		if s, ok := db[key]; ok {
			return []byte(s), nil
		}
		return nil, fmt.Errorf("key not found from db")
	}))

	addr := "localhost:9999"
	pool := gocache.NewHTTPPool(addr)
	log.Fatal(http.ListenAndServe(addr, pool))

}

func TestSearch() {
	a := []int{1, 3, 5, 10}
	f := 0
	ind := sort.Search(len(a), func(i int) bool {
		return a[i] >= f
	})
	fmt.Println(ind)
}


func main() {
	//TestGetter()
	//testGroup()
	//testHttp()
	TestSearch()
}
