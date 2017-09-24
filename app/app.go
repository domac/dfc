package app

import (
	"errors"
	c "github.com/domac/dfc/cache"
	"github.com/domac/dfc/store"
	"log"
	"runtime/debug"
)

var started bool

var ErrNullCacheServer = errors.New("cache server was null")

//全局默认缓存服务
var DefaultCacheServer *DFCServer
var DefaultPeerRoundRobin *SessionPeers
var DefaultResourceDB *store.ResourceDB

func GetCacheServer() (*DFCServer, error) {
	if DefaultCacheServer == nil {
		return nil, ErrNullCacheServer
	}
	return DefaultCacheServer, nil
}

func GetResourceDB() (*store.ResourceDB, error) {
	return DefaultResourceDB, nil
}

//总服务开关
func Startup(cfg *AppConfig) (err error) {
	if started {
		return
	}

	DefaultCacheServer = NewDFCServer(cfg)

	//启动会话服务
	sessionServer, err := NewSessionServer(cfg)
	if err == nil {
		sessionServer.Start()
	}

	DefaultResourceDB, err = store.OpenResourceDB(cfg.Local_store_path)
	if err != nil {
		log.Println(err.Error())
	}

	peerInfos := cfg.Peer
	DefaultPeerRoundRobin, _ = NewSessionPeers(peerInfos)

	started = true
	return
}

//停止服务
func Shutdown(i interface{}) {
	println()
	if DefaultResourceDB != nil {
		DefaultResourceDB.Close()
	}
	log.Println("application ready to shutdown")
}

//处理缓存请求服务
type DFCServer struct {
	cache *c.Cache //free cache
	cfg   *AppConfig
}

func NewDFCServer(cfg *AppConfig) *DFCServer {
	log.Printf("dfc max cache object size: %d | expired seconds: %d\n",
		cfg.Cache_max_size, cfg.Cache_ttl)
	debug.SetGCPercent(20)
	return &DFCServer{
		cache: c.NewCache(cfg.Cache_max_size),
		cfg:   cfg,
	}
}

func (d *DFCServer) GetConfig() *AppConfig {
	return d.cfg
}

func (d *DFCServer) Cache() *c.Cache {
	return d.cache
}

func (d *DFCServer) Set(key string, value []byte) (err error) {
	err = d.cache.Set([]byte(key), value, d.cfg.Cache_ttl)
	return
}

func (d *DFCServer) Get(key string) (value []byte, err error) {
	value, err = d.cache.Get([]byte(key))
	return
}
