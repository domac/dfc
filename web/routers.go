package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//载入路由表
func loadRouters() (s *http.Server, err error) {

	//图片服务
	imageHandler := &ImageHandler{}

	r := mux.NewRouter()
	v1Subrouter := r.PathPrefix("/v1").Subrouter()
	ih := NewImageHandler(imageHandler.Resize)
	v1Subrouter.Handle("/{image}", ih).Methods("GET") //GET的响应处理

	//接入http server
	s = &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:10200",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s, nil
}
