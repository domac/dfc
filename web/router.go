package web

import (
	"github.com/gorilla/mux"
)

//载入路由表
func loadRouter() (r *mux.Router, err error) {
	//图片服务
	imageHandler := &ImageHandler{}
	r = mux.NewRouter()
	v1Subrouter := r.PathPrefix("/v1").Subrouter()
	ih := NewImageHandler(imageHandler.Resize)
	v1Subrouter.Handle("/image.do", ih).Methods("GET") //GET的响应处理
	return r, nil
}
