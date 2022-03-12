# go-cache(分布式缓存)

go mod init go-cache

## 使用方式
demo: 运行main.go

引入使用：
lru.New(系统容量, 回调函数)

测试tag

## day4  实现一致性HASH模块

1. 在多个存储节点时候，需要实现一致性hash
2. 一致性hash 将节点映射到多个节点中间，如果节点数较少，通过虚拟节点方式
   增加节点，使得数据更加分散
3. 在获取数据时候，也是得到节点hash值，扫描数据
