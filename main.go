package main

import (
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/web"
	"log"
)

func main() {

	println(app.Version)

	//启用应用服务
	if err := app.Startup(); err != nil {
		log.Fatal(err)
		return
	}

	//启用服务端程序
	httpServer, err := web.InitServer(":10200")
	if err != nil {
		log.Fatal(err)
		return
	}

	//启用web API
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err.Error())
		}
	}()

	// 注册退出事件, 回调函数用于清理工作
	app.On(app.EXIT, app.Shutdown)
	// 监听退出信号
	app.Wait()
	app.Emit(app.EXIT, nil)
	log.Println("application is exit !")
}
