package server

import "github.com/gin-gonic/gin"

// BaseRouter 基础路由结构
type BaseRouter struct {
	Group *gin.RouterGroup
}

// RegisterGroup 注册路由组的基础方法
func RegisterGroup(r *gin.Engine, path string) *BaseRouter {
	return &BaseRouter{
		Group: r.Group(path),
	}
}

// 路由注册接口
type IRouter interface {
	Register(r *gin.Engine)
}
