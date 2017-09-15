package web

import (
	"github.com/gorilla/mux"
)

//载入路由表
func loadRouter() (r *mux.Router, err error) {
	//缓存服务
	cacheHandler := &CacheHandler{}
	r = mux.NewRouter()
	v1Subrouter := r.PathPrefix("/v1").Subrouter()
	ih := NewCacheHandler(cacheHandler.Cache)
	v1Subrouter.Handle("/cache.do", ih).Methods("GET") //GET的响应处理
	return r, nil
}
