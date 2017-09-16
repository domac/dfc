package app

import (
	c "github.com/domac/dfc/cache"
	"log"
	"runtime/debug"
)

var started bool

//global cache server
var DefaultCacheServer *DFCServer

func GetCacheServer() *DFCServer {
	if DefaultCacheServer == nil {
		log.Println("restruct a new cache server")
		DefaultCacheServer = NewDFCServer(1024 * 1024 * 1024)
	}
	return DefaultCacheServer
}

//Application Run
//Just run once
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

func Shutdown(i interface{}) {
	println()
	log.Println("application ready to shutdown")
}

type DFCServer struct {
	cache *c.Cache
}

func NewDFCServer(cacheSize int) *DFCServer {
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
