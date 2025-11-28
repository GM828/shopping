package service

import (
	"shopping/product_service/model"
	"shopping/product_service/repository"
)

type ProductService interface {
	// 根据商品Id查询商品
	GetProductById(request *model.ProductMO) (*model.ProductMO, error)
	// 根据商品Id查询商品库存
	GetStockByProductId(m *model.ProductStockMO) (*model.ProductStockMO, error)
}

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		productRepository: productRepository,
	}
}

func (s ProductServiceImpl) GetProductById(mo *model.ProductMO) (*model.ProductMO, error) {
	return s.productRepository.GetProductById(mo)
}

func (s ProductServiceImpl) GetStockByProductId(mo *model.ProductStockMO) (*model.ProductStockMO, error) {
	return s.productRepository.GetStockByProductId(mo)
}
