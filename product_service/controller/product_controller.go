package controller

import (
	"github.com/gin-gonic/gin"
	"shopping/product_service/dto"
	"shopping/product_service/manager"
)

type ProductController struct {
	ProductManager manager.ProductManager
}

func (c *ProductController) GetProductById(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.ProductRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	response, err := c.ProductManager.GetProductById(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, response)
}

func (c *ProductController) CreateIndex(ctx *gin.Context) {
	err := c.ProductManager.CreateIndex()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "索引创建成功"})
}

func (c *ProductController) GetProductsByName(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.ProductRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	response, err := c.ProductManager.GetProductsByName(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, response)
}
