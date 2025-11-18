package user

import (
	"github.com/gin-gonic/gin"
	"shopping/server"
)

type Route struct {
	server.BaseRouter
}

func (r *Route) Register(engine *gin.Engine) {
	userController, err := InitializeUserController()
	if err != nil {
		panic(err)
	}
	userGroup := engine.Group("/client/user")
	userGroup.POST("/login", userController.Login)
	userGroup.POST("/register", userController.Register)
	userGroup.POST("/updatePassword", userController.UpdatePassword)
	userGroup.POST("/updateUserInfo", userController.UpdateUserInfo)
}
