package main

import (
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/web"
	"log"
	"os"
)

func main() {
	println(app.Version)
	httpServer, err := web.InitServer()
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	//开启web服务
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err.Error())
		}
	}()

	// 注册退出事件, 回调函数用于清理工作
	app.On(app.EXIT, func(i interface{}) {
		log.Println("app is ready to close")
	})
	// 监听退出信号
	app.Wait()
	app.Emit(app.EXIT, nil)
	println("exit success")
}
