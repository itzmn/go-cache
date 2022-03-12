package gocache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// http服务

const (
	project         = "/gocache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string                 // 记录项目 ip 端口
	basePath    string                 // 记录项目名称
	mu          sync.Mutex             // 读写锁
	peers       *Map                   // 一致性hash的map
	httpGetters map[string]*httpGetter //  存储每个远程节点 对于的 getter
}

// 构建 http 服务端

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: project,
		//httpGetters: make(map[string]*httpGetter),
		//peers: NewMap(defaultReplicas, nil),
	}
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf(format, v)
}

// 实现http 的接口
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool Serve error path:" + r.URL.Path)
	}

	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)

	// 请求 服务不对
	if group == nil {
		http.Error(w, "no such group:"+groupName, http.StatusBadRequest)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())

}

// Set 针对服务端，初始化 一致性hash模块
func (h *HTTPPool) Set(peers ...string) {

	h.mu.Lock()
	// 设置 一致性hash
	h.peers = NewMap(defaultReplicas, nil)
	h.mu.Unlock()

	// 给一致性hash 模块 增加机器
	h.peers.Add(peers...)

	h.httpGetters = make(map[string]*httpGetter, len(peers))
	// 给每一个 远程节点分配 getter客户端
	for _, peer := range peers {

		h.httpGetters[peer] = &httpGetter{
			baseUrl: peer + h.basePath,
		}
	}

}

// PickPeer 实现PeerPicker 接口， 根据传入的key 获取对应的 获取缓存的getter 客户端
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {

	h.mu.Lock()
	defer h.mu.Unlock()
	// 根据key 得到真实的存储节点, 不是自己, 也不是空
	if get := h.peers.Get(key); get != "" && get != h.self{
		h.Log("Pick peer: %v", get)
		return h.httpGetters[get], true
	}
	return nil, false
}


var _ PeerPicker = (*HTTPPool)(nil)


// =========== 客户端

type httpGetter struct {
	baseUrl string // 前缀路径
}

// Get 实现PeerGetter接口 从group 获取缓存的方法，  作为一个 客户端
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	sprintf := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	res, err := http.Get(sprintf)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return:%v", res.Status)
	}

	all, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("read response body err: %v", err)
	}

	return all, nil
}

var _ PeerGetter = (*httpGetter)(nil)
