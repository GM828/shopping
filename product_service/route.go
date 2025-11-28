package product_service

import (
	"github.com/gin-gonic/gin"
	"shopping/server"
)

type Route struct {
	server.BaseRouter
}

func (r *Route) Register(engine *gin.Engine) {
	productController, err := InitializeProductController()
	if err != nil {
		panic(err)
	}
	productGroup := engine.Group("/product")
	productGroup.POST("/getProduct", productController.GetProductById)
}
