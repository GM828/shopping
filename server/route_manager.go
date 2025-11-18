package server

import "github.com/gin-gonic/gin"

// RouterManager 路由管理器
type RouterManager struct {
	engine  *gin.Engine
	routers []IRouter
}

// NewRouterManager 创建路由管理器
func NewRouterManager(engine *gin.Engine) *RouterManager {
	return &RouterManager{
		engine:  engine,
		routers: make([]IRouter, 0),
	}
}

// Register 注册路由
func (rm *RouterManager) Register(router IRouter) {

	rm.routers = append(rm.routers, router)
}

// Init 初始化所有路由
func (rm *RouterManager) Init() {
	for _, router := range rm.routers {
		router.Register(rm.engine)
	}
}
