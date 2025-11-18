//go:build wireinject
// +build wireinject

package user

import (
	"github.com/google/wire"
	"shopping/thrift_client/user/controller"
)

func InitializeUserController() (*controller.UserController, error) {
	wire.Build(
		// Controller
		wire.Struct(new(controller.UserController), "*"),
	)
	return nil, nil
}
