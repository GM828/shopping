package handler

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"shopping/product_service/dto"
	"shopping/product_service/manager"
	product "shopping/product_service/server/kitex_gen/product"
)

// ProductServiceImpl implements the last service interface defined in the IDL.
type ProductServiceImpl struct {
	ProductManager manager.ProductManager
}

// GetProductById implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProductById(ctx context.Context, request *product.ProductRequest) (resp *product.ProductResponse, err error) {
	respDto, err := s.ProductManager.GetProductById(&dto.ProductRequestDTO{
		Id: request.Id,
	})
	if err != nil {
		// 这里是业务异常，根据Kitex的语法，需要自定义异常
		err = kerrors.NewBizStatusError(500, err.Error())
		return nil, err
	}
	resp = &product.ProductResponse{
		Id:          respDto.Id,
		ProductName: respDto.ProductName,
		Category:    respDto.Category,
		Description: respDto.Description,
		Price:       respDto.Price,
		Status:      respDto.Status,
		CreateTime:  respDto.CreateTime,
		StockNum:    respDto.StockNum,
	}
	return resp, nil
}
