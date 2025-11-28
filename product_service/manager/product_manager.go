package manager

import (
	"errors"
	"shopping/product_service/dto"
	"shopping/product_service/model"
	"shopping/product_service/service"
	"shopping/util"
)

type ProductManager interface {
	GetProductById(request *dto.ProductRequestDTO) (*dto.ProductResponseDTO, error)
}
type ProductManagerImpl struct {
	productService service.ProductService
}

func NewProductManager(productService service.ProductService) ProductManager {
	return &ProductManagerImpl{
		productService: productService,
	}
}

func (m ProductManagerImpl) GetProductById(request *dto.ProductRequestDTO) (*dto.ProductResponseDTO, error) {
	// 参数校验
	if request.Id == nil {
		return nil, errors.New("product id is required")
	}
	// 查询
	mo, err := m.productService.GetProductById(&model.ProductMO{
		Id: request.Id,
	})
	if err != nil {
		return nil, err
	}
	if mo == nil {
		return nil, errors.New("product is not exist, id: " + util.Int64ToStr(*request.Id))
	}
	// 查询库存
	stockMo, err := m.productService.GetStockByProductId(&model.ProductStockMO{
		ProductId: mo.Id,
	})
	if err != nil {
		return nil, err
	}
	// 封装返回数据
	createTime := util.DateUtil.FormatDateByCustomLayout(mo.CreateTime, util.DateLayout.YYYY_MM_DD_HH_MM_SS)
	response := &dto.ProductResponseDTO{
		Id:          mo.Id,
		ProductName: mo.ProductName,
		Category:    mo.Category,
		Price:       mo.Price,
		Description: mo.Description,
		Status:      mo.Status,
		CreateTime:  &createTime,
		StockNum:    stockMo.StockNum,
	}
	return response, nil
}
