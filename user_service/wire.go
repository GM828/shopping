//go:build wireinject
// +build wireinject

package user_service

import (
	"github.com/google/wire"
	"shopping/user_service/controller"
	"shopping/user_service/manager"
	"shopping/user_service/po"
	"shopping/user_service/repository"
	"shopping/user_service/server/handler"
	"shopping/user_service/service"
)

func InitializeUserController() (*controller.UserController, error) {
	wire.Build(
		// 新增：PO 的提供者（必须放在 Repository 前面，因为 Repository 依赖 PO）
		po.NewUserInfoPO,
		po.NewUserLoginPO,

		// Repository
		repository.NewUserRepository,
		// Service
		service.NewUserService,
		// Manager
		manager.NewUserManager,
		// Controller
		wire.Struct(new(controller.UserController), "*"),
	)
	return nil, nil
}

func InitializeUserHandler() (*handler.UserHandler, error) {
	wire.Build(
		// 新增：PO 的提供者（必须放在 Repository 前面，因为 Repository 依赖 PO）
		po.NewUserInfoPO,
		po.NewUserLoginPO,

		// Repository
		repository.NewUserRepository,
		// Service
		service.NewUserService,
		// Manager
		manager.NewUserManager,
		// Controller
		wire.Struct(new(handler.UserHandler), "*"),
	)
	return nil, nil
}
