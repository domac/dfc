package app

import (
	"log"
)

var started bool

//服务初始化
func Startup() (err error) {
	if started {
		return
	}

	//1.初始化配置
	//2.初始化db层
	//3.初始化业务功能

	ImgCli = InitImageClient()

	started = true
	return
}

//服务关闭退出
func Shutdown(i interface{}) {
	println()
	log.Println("application ready to shutdown")
}
