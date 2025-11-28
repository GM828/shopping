package main

import (
	"github.com/gin-gonic/gin"
	"shopping/database"
	"shopping/product_service"
	"shopping/server"
	"shopping/user_service"
)

func main() {
	// 初始化gin
	r := gin.Default()
	// 创建路由管理器
	routerManager := server.NewRouterManager(r)
	// 注册路由
	RegisterRoutes(routerManager)
	// 初始化数据库
	database.Init()

	err := r.Run(":9090")
	if err != nil {
		return
	}
}

func RegisterRoutes(routerManager *server.RouterManager) {
	routerManager.Register(&user_service.Route{})
	routerManager.Register(&product_service.Route{})
	// 初始化所有路由
	routerManager.Init()
}
