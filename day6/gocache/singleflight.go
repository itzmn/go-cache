package gocache

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type FGroup struct {
	mu sync.Mutex       // 调用类， 并发锁
	m  map[string]*call // 保存不同key 的调用
}

// DO 接收key 和func， 同个key func 只会执行一次
func (f *FGroup) DO(key string, fn func() (interface{}, error)) (interface{}, error) {

	if c, ok := f.m[key]; ok {
		// 该key的调用函数正在执行
		c.wg.Wait()
		return c.val, nil
	}

	c := new(call)

	// 发起请求之前需要 枷锁
	c.wg.Add(1)
	// 表明 key 已经有人在处理
	f.m[key] = c
	c.val, c.err = fn()
	c.wg.Done()
	delete(f.m, key)

	return c.val, c.err
}
