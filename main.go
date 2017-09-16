package main

import (
	"github.com/domac/dfc/app"
	"github.com/domac/dfc/web"
	"log"
	_ "net/http/pprof"
)

//prof command:
//go tool pprof --seconds 50 http://localhost:10200/debug/pprof/profile
func main() {

	println(app.Version)

	//start up app server
	if err := app.Startup(); err != nil {
		log.Fatal(err)
		return
	}

	//open http server
	httpServer, err := web.InitServer(":10200")
	if err != nil {
		log.Fatal(err)
		return
	}

	//open web api
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err.Error())
		}
	}()

	//register some event when user press `Ctrl + C`
	app.On(app.EXIT, app.Shutdown)
	app.Wait()
	app.Emit(app.EXIT, nil)
	log.Println("dfc is exit !")
}
