package main

import (
	"red_packet/initialize"
	"red_packet/router"

	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化App
	app := initialize.NewApp()
	// 连接数据库，缓存预热等操作
	app.Run()
	// 初始化路由
	r := router.InitRouter()
	logrus.Error(r.Listen(":8080"))
}
