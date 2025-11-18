package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shopping/database"
	"shopping/server"
	"shopping/thrift_client/clients"
	"shopping/thrift_client/thrift_manager"
	"shopping/thrift_client/user"
	"syscall"
	"time"
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

	// 初始化Thrift管理器
	thriftManager := thrift_manager.NewThriftManager()
	// 注册Thrift服务
	RegisterThriftServices(thriftManager)

	// 启动HTTP服务
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		log.Println("服务启动成功，端口: 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务启动失败: %v", err)
		}
	}()

	// 监听退出信号，批量关闭资源
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("开始关闭服务...")

	// 关闭HTTP服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP服务关闭异常: %v", err)
	}

	// 批量关闭所有Thrift连接（通过管理器）
	if err := thriftManager.Close(); err != nil {
		log.Printf("Thrift连接关闭异常: %v", err)
	}

	log.Println("服务已退出")
}

func RegisterRoutes(routerManager *server.RouterManager) {
	routerManager.Register(&user.Route{})
	// 初始化所有路由
	routerManager.Init()
}

func RegisterThriftServices(thriftManager *thrift_manager.ThriftManager) {
	clients.RegisterUserService(thriftManager, "localhost:9090", 10)
	// 初始化所有Thrift服务
	err := thriftManager.Init()
	if err != nil {
		panic(err)
	}
}
