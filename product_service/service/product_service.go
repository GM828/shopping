package service

import (
	"shopping/product_service/model"
	"shopping/product_service/repository"
)

type ProductService interface {
	MySQLService
	EsService
}

// MySQL相关业务操作
type MySQLService interface {
	// 根据商品Id查询商品
	GetProductById(request *model.ProductMO) (*model.ProductMO, error)
	// 根据商品Id查询商品库存
	GetStockByProductId(m *model.ProductStockMO) (*model.ProductStockMO, error)
}

// Es相关业务操作
type EsService interface {
	// 创建索引
	CreateIndex(mo *model.ProductFullMO) error
	// 模糊查询文档
	SearchDocs(query string) ([]*model.ProductFullMO, error)
}

type ProductServiceImpl struct {
	productRepository repository.ProductRepository
	esRepository      repository.EsRepository
}

func NewProductService(productRepository repository.ProductRepository, esRepository repository.EsRepository) ProductService {
	return &ProductServiceImpl{
		productRepository: productRepository,
		esRepository:      esRepository,
	}
}

func (s ProductServiceImpl) GetProductById(mo *model.ProductMO) (*model.ProductMO, error) {
	return s.productRepository.GetProductById(mo)
}

func (s ProductServiceImpl) GetStockByProductId(mo *model.ProductStockMO) (*model.ProductStockMO, error) {
	return s.productRepository.GetStockByProductId(mo)
}

func (s ProductServiceImpl) CreateIndex(mo *model.ProductFullMO) error {
	return s.esRepository.CreateIndex(mo.Mapping())
}

func (s ProductServiceImpl) SearchDocs(query string) ([]*model.ProductFullMO, error) {
	return s.esRepository.SearchDocs(query)
}
