//go:build wireinject
// +build wireinject

package product_service

import (
	"github.com/google/wire"
	"shopping/product_service/controller"
	"shopping/product_service/manager"
	"shopping/product_service/po"
	"shopping/product_service/repository"
	"shopping/product_service/server/handler"
	"shopping/product_service/service"
)

func InitializeProductController() (*controller.ProductController, error) {
	wire.Build(
		// 新增：PO 的提供者（必须放在 Repository 前面，因为 Repository 依赖 PO）
		po.NewProductPO,
		po.NewProductStockPO,

		// Repository
		repository.NewProductRepository,
		repository.NewEsRepository,
		// Service
		service.NewProductService,
		// Manager
		manager.NewProductManager,
		// Controller
		wire.Struct(new(controller.ProductController), "*"),
	)
	return nil, nil
}

func InitializeProductHandler() (*handler.ProductServiceImpl, error) {
	wire.Build(
		// 新增：PO 的提供者（必须放在 Repository 前面，因为 Repository 依赖 PO）
		po.NewProductPO,
		po.NewProductStockPO,

		// Repository
		repository.NewProductRepository,
		repository.NewEsRepository,
		// Service
		service.NewProductService,
		// Manager
		manager.NewProductManager,
		// Controller
		wire.Struct(new(handler.ProductServiceImpl), "*"),
	)
	return nil, nil
}
