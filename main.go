package main

import (
	"flag"
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/web"
	"log"
	_ "net/http/pprof"
)

var (
	config = flag.String("config", "./conf/base.conf", "set the config file path")
)

//prof command:
//go tool pprof --seconds 50 http://localhost:10200/debug/pprof/profile
func main() {

	println(app.Version)
	flag.Parse()

	log.Printf("app config file : %s\n", *config)
	cfg, err := app.LoadConfig(*config)
	if err != nil {
		log.Fatal(err)
		return
	}

	//start up app server
	if err := app.Startup(cfg); err != nil {
		log.Fatal(err)
		return
	}

	//open http server
	log.Printf("app is listening %s", cfg.Http_address)
	httpServer, err := web.InitServer(cfg.Http_address)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err.Error())
		}
	}()

	//注册退出事件
	app.On(app.EXIT, app.Shutdown)
	app.Wait()
	app.Emit(app.EXIT, nil)
	log.Println("dfc is exit now !")
}
