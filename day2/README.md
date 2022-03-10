# go-cache(分布式缓存)

go mod init go-cache

## 使用方式
demo: 运行main.go

引入使用：
lru.New(系统容量, 回调函数)

测试tag

## day2  对缓存模块进行封装

1、封装 缓存值类型 ByteView
2、对缓存并发模块进行封装 cache
3、增加缓存获取不到时候的回调函数
4、封装Group模块，隔离底层，用户操作类
