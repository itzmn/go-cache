package gocache

// 抽象 PeerPicker

// PeerGetter 从group 中获取缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

type PeerPicker interface {
	// PickPeer 根据传入的key 得到对于的缓存 获取类
	PickPeer(key string)(peer PeerGetter, ok bool)
}



