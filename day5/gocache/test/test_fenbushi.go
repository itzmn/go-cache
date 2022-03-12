package main

import (
	"fmt"
	"gocache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup(group string) *gocache.Group {
	return gocache.NewGroup(group, 2<<5, gocache.GetterFunc(func(key string) ([]byte, error) {
		log.Printf("search %s from db", key)
		if s, ok := db[key]; ok {
			return []byte(s), nil
		}
		return nil, fmt.Errorf("key not found from db")
	}))

}

func startCacheServer(addr string, addrs []string, group *gocache.Group) {
	peers := gocache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeerPicker(peers)
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *gocache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(view.ByteSlice())
	}))

	log.Println("api server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {

	//var port int
	//var api bool
	//flag.IntVar(&port, "port", 8081, "服务端口")
	//flag.BoolVar(&api, "api", false, "是否是api服务器")
	//flag.Parse()

	port := 8001
	api := false

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup("scores")
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], addrs, gee)

}
