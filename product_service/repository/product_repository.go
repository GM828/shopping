package repository

import (
	"shopping/product_service/model"
	"shopping/product_service/po"
)

type ProductRepository interface {
	// 根据商品Id查询商品
	GetProductById(request *model.ProductMO) (*model.ProductMO, error)
	// 根据商品Id查询商品库存
	GetStockByProductId(m *model.ProductStockMO) (*model.ProductStockMO, error)
}

type ProductRepositoryImpl struct {
	productPO      *po.ProductPO
	productStockPO *po.ProductStockPO
}

func NewProductRepository(productPO *po.ProductPO, productStockPO *po.ProductStockPO) ProductRepository {
	return &ProductRepositoryImpl{
		productPO:      productPO,
		productStockPO: productStockPO,
	}
}

func (r ProductRepositoryImpl) GetProductById(mo *model.ProductMO) (*model.ProductMO, error) {
	return r.productPO.FindOne(mo)
}

func (r ProductRepositoryImpl) GetStockByProductId(mo *model.ProductStockMO) (*model.ProductStockMO, error) {
	return r.productStockPO.FindOne(mo)
}
