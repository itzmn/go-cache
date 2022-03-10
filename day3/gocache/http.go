package gocache

import (
	"log"
	"net/http"
	"strings"
)

// http服务

const project = "/gocache/"

type HTTPPool struct {
	self     string // 记录项目 ip 端口
	basePath string // 记录项目名称
}


// 构建 http 服务端

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: project,
	}
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf(format, v)
}

// 实现http 的接口
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request)  {

	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool Serve error path:" + r.URL.Path )
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
		http.Error(w, "no such group:" + groupName, http.StatusBadRequest)
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