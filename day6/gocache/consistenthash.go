package gocache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 分布式 一致性hash模块

// Hash 定义 hash 杉树类型
type Hash func(data []byte) uint32

// Map 定义容器，存储所有的数据
type Map struct {
	hash     Hash           // 指定hash函数
	keys     []int          // 存储所有的hash 位置，排序的
	replicas int            // 指定机器的虚拟倍数
	hashMap  map[int]string //存储虚拟节点和 真实机器的映射
}

// NewMap 根据传入的虚拟个数 和 hash函数 创建容器
func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 根据传入的多个真实机器， 组件存储容器
func (m *Map) Add(machines ...string) {
	for _, machine := range machines {
		for i := 0; i < m.replicas; i++ {
			// 根据 副本索引和 机器名称 得到 节点的hash值
			hash := int(m.hash([]byte(strconv.Itoa(i) + machine)))
			m.keys = append(m.keys, hash)
			// 保存环上的hash 值和 真实机器的关系
			m.hashMap[hash] = machine
		}
	}
	// 对hash 列表 排序
	sort.Ints(m.keys)
}

// Get 根据需要存储的key 选择真实存储的节点
func (m *Map) Get(key string) string {

	if len(m.keys) == 0 {
		return ""
	}

	// 得到key 的hash值
	hash := int(m.hash([]byte(key)))

	// 查找hash值 对应节点位置，
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 获取真实机器
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
