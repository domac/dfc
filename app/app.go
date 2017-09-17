package app

import (
	"errors"
	c "github.com/domac/dfc/cache"
	"log"
	"runtime/debug"
)

var started bool

var ErrNullCacheServer = errors.New("cache server was null")

//全局默认缓存服务
var DefaultCacheServer *DFCServer

func GetCacheServer() (*DFCServer, error) {
	if DefaultCacheServer == nil {
		return nil, ErrNullCacheServer
	}
	return DefaultCacheServer, nil
}

//总服务开关
func Startup(configPath string) (err error) {
	if started {
		return
	}

	log.Printf("config file : %s\n", configPath)

	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	DefaultCacheServer = NewDFCServer(cfg.Cache_max_size)

	sessionServer := NewSessionServer(cfg.Max_reqs_per_second)
	sessionServer.Start()

	started = true
	return
}

func Shutdown(i interface{}) {
	println()
	log.Println("application ready to shutdown")
}

//处理缓存请求服务
type DFCServer struct {
	cache *c.Cache //free cache
}

func NewDFCServer(cacheSize int) *DFCServer {
	log.Printf("dfc max cache object size: %d\n", cacheSize)
	debug.SetGCPercent(20)
	return &DFCServer{
		cache: c.NewCache(cacheSize),
	}
}

func (d *DFCServer) Cache() *c.Cache {
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
