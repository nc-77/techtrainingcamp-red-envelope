package main

import (
	"red_envelope/router"
	"red_envelope/service"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化App
	app := service.GetApp()
	// 连接数据库，缓存预热等操作
	app.Run()
	// 初始化路由
	r := router.InitRouter()
	logrus.Error(r.Listen(":8080"))
}
