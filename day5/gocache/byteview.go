package gocache

// 缓存数据的抽象

type ByteView struct {
	b []byte
}

// 实现 Value接口方法 作为底层数据 一种

func (b ByteView) Len() int {
	return len(b.b)
}

func cloneBytes(data []byte) []byte  {
	bytes := make([]byte, len(data))
	copy(bytes, data)
	return bytes
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}