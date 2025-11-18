package server

import (
	"github.com/gin-gonic/gin"
)

func InitGin(appName string) *gin.Engine {
	/*
		判断生产环境
		if config.IsPrd() {
			gin.SetMode(gin.ReleaseMode)
		} else {
			gin.SetMode(gin.DebugMode)
		}*/

	r := gin.New()
	/*r.Use(cors.Default())
	r.Use(otelgin.Middleware(appName))
	// 添加gin的otel中间件
	r.Use(otelgin.Middleware(config.Instance.GetString("appName")))
	r.Use(func(c *gin.Context) {
		otel.SetRequestContext(c.Request.Context())
		defer otel.RemoveRequestContext()
		c.Next()
	})
	//r.Use(gin.Recovery())
	r.Use(recoverFn())
	// Add global middlewares
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//全局异常捕获
	r.Use(middleware.GloableRecoverFn())
	// 注册全局异常处理器
	r.Use(middleware.ExceptionResponseHandler())
	// 请求响应日志
	r.Use(RequestResponseLogger())*/

	return r
}
