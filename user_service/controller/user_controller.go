package controller

import (
	"github.com/gin-gonic/gin"
	"shopping/user_service/dto"
	"shopping/user_service/manager"
)

type UserController struct {
	UserManager manager.UserManager
}

func (c *UserController) Login(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserLoginRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	response, err := c.UserManager.UserLogin(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, response)
}

func (c *UserController) Register(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserRegisterRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	err := c.UserManager.UserRegister(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, gin.H{"message": "注册成功"})
}

func (c *UserController) UpdatePassword(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserUpdateLoginRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	err := c.UserManager.UpdatePassword(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, gin.H{"message": "密码更新成功"})
}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserUpdateInfoRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 调用manager层方法
	err := c.UserManager.UpdateUserInfo(request)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 返回响应
	ctx.JSON(200, gin.H{"message": "用户信息更新成功"})
}
