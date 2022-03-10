# go-cache(分布式缓存)

go mod init go-cache

## 使用方式
demo: 运行main.go

引入使用：
lru.New(系统容量, 回调函数)

测试tag

## day3  在HTTP服务中使用缓存

1、创建缓存模块
2、创建HTTP服务，在处理函数中使用缓存模块
