package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"shopping/thrift_client/clients"
	"shopping/thrift_client/user/dto"
	"shopping/thrift_def/gen-go/user"
)

type UserController struct {
}

func (c *UserController) Login(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserLoginRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取UserService的客户端
	userClient, wrapped, err := clients.BorrowUserClient()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "获取 user 客户端失败: " + err.Error()})
		return
	}
	// 日志：显示当前获取到的连接池对象标识
	log.Printf("[连接池] 获取连接 | 对象地址: %p | Transport地址: %p | 是否打开: %v\n",
		wrapped,
		wrapped.Transport,
		wrapped.Transport.IsOpen())

	// 关键：归还客户端连接
	defer clients.ReturnUserClient(wrapped)

	req := &user.UserLoginRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	resp, err := userClient.Login(context.Background(), req)
	if err != nil {
		// 处理异常
		if ex, ok := err.(*user.UserException); ok {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("登录失败: %s", ex.Message)})
		} else {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("调用错误: %v", err)})
		}
		return
	}
	// 返回响应
	ctx.JSON(200, resp)
}

func (c *UserController) Register(ctx *gin.Context) {
	// 接收请求参数
	request := &dto.UserRegisterRequestDTO{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 获取UserService的客户端
	userClient, wrapped, err := clients.BorrowUserClient()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "获取 user 客户端失败: " + err.Error()})
		return
	}
	// 关键：归还客户端连接
	defer clients.ReturnUserClient(wrapped)

	req := &user.UserRegisterRequest{
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
		RealName: request.RealName,
		Phone:    request.Phone,
		Gender:   request.Gender,
		Birthday: request.Birthday,
	}
	_, err = userClient.Register(context.Background(), req)
	if err != nil {
		// 处理异常
		if ex, ok := err.(*user.UserException); ok {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("注册失败: %s", ex.Message)})
		} else {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("调用错误: %v", err)})
		}
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

	// 获取UserService的客户端
	userClient, wrapped, err := clients.BorrowUserClient()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "获取 user 客户端失败: " + err.Error()})
		return
	}
	// 关键：归还客户端连接
	defer clients.ReturnUserClient(wrapped)

	req := &user.UserUpdateLoginRequest{
		UserName:     request.UserName,
		Email:        request.Email,
		Password:     request.Password,
		NewPassword_: request.NewPassword,
	}
	_, err = userClient.UpdatePassword(context.Background(), req)
	if err != nil {
		// 处理异常
		if ex, ok := err.(*user.UserException); ok {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("更新密码失败: %s", ex.Message)})
		} else {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("调用错误: %v", err)})
		}
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

	// 获取UserService的客户端
	userClient, wrapped, err := clients.BorrowUserClient()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "获取 user 客户端失败: " + err.Error()})
		return
	}
	// 关键：归还客户端连接
	defer clients.ReturnUserClient(wrapped)

	req := &user.UserUpdateInfoRequest{
		UserInfoId: request.UserInfoId,
		RealName:   request.RealName,
		Phone:      request.Phone,
		Gender:     request.Gender,
		Birthday:   request.Birthday,
	}
	_, err = userClient.UpdateUserInfo(context.Background(), req)
	if err != nil {
		// 处理异常
		if ex, ok := err.(*user.UserException); ok {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("更新用户信息失败: %s", ex.Message)})
		} else {
			ctx.JSON(500, gin.H{"error": fmt.Sprintf("调用错误: %v", err)})
		}
		return
	}

	// 返回响应
	ctx.JSON(200, gin.H{"message": "用户信息更新成功"})
}
