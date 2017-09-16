package app

import (
	"github.com/coocood/freecache"
	"log"
	"runtime/debug"
)

var started bool

var DefaultCacheServer *DFCServer

func GetCacheServer() *DFCServer {
	if DefaultCacheServer == nil {
		log.Println("restruct a new cache server")
		DefaultCacheServer = NewDFCServer(1024 * 1024 * 1024)
	}
	return DefaultCacheServer
}

//服务初始化
func Startup() (err error) {
	if started {
		return
	}
	DefaultCacheServer = NewDFCServer(1024 * 1024 * 1024)

	sessionServer := NewSessionServer(8000)
	sessionServer.Start()

	started = true
	return
}

//服务关闭退出
func Shutdown(i interface{}) {
	println()
	log.Println("application ready to shutdown")
}

//DFC 服务
type DFCServer struct {
	cache *freecache.Cache
}

func NewDFCServer(cacheSize int) *DFCServer {
	debug.SetGCPercent(20)
	return &DFCServer{
		cache: freecache.NewCache(cacheSize),
	}
}

//获取缓存存储
func (d *DFCServer) Cache() *freecache.Cache {
	return d.cache
}

func (d *DFCServer) Set(key string, value []byte, expireSeconds int) (err error) {
	err = d.cache.Set([]byte(key), value, expireSeconds)
	return
}

func (d *DFCServer) Get(key string) (value []byte, err error) {
	value, err = d.cache.Get([]byte(key))
	return
}
