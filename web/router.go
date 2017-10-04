package web

import (
	"github.com/gorilla/mux"
)

//load router table
func loadRouter() (r *mux.Router, err error) {
	cacheHandler := &CacheHandler{}
	r = mux.NewRouter()
	v1Subrouter := r.PathPrefix("/v1").Subrouter()
	ih := NewCacheHandler(cacheHandler.Cache)
	v1Subrouter.Handle("/cache.do", ih).Methods("GET") //GET Method

	//缓存状态信息
	cs := NewCacheHandler(cacheHandler.CacheStats)
	v1Subrouter.Handle("/stats.do", cs).Methods("GET") //GET Method
	return r, nil
}
