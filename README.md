# DFC

简单的文件缓存服务

## 介绍

日常我们经常有大量请求服务器文件的场景，这些请求万一非常频繁访问机器的本地磁盘, 这样会对服务器IO产生很大的影响，严重的情况会影响其他服务的使用。于是针对这种高风险的情况，我开发了这个工具让每次的文件获取请求，都尽量能通过命中 in-momory和local kv cache,而不是每次都访问物理磁盘，这样会提高请求吞吐量和性能。

目前DFC支持简单的水平扩容（parent模式），通过部署多个dfc，并配置上下层级关系，就能让缓存数据得到某程度上的“共享”，提高请求命中率。

## 使用方式

```
$ go get github.com/domac/dfc
```

进入程序目录，执行相关命令

```
$ make && cd release

$ ./dfc -conf=base.conf
```

## 功能API

### 数据请求

示例：
```
http://localhost:10200/v1/cache.do?url=/your/file/path
```

### 缓存统计信息

示例：

```
http://localhost:10200/v1/stats.do
```