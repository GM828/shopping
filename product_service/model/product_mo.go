package model

import "time"

// 商品核心信息
type ProductMO struct {
	Id          *int64     `json:"id"`          // 商品ID
	ProductName *string    `json:"productName"` // 商品名称
	Category    *string    `json:"category"`    // 商品分类（直接存储名称）
	Price       *float64   `json:"price"`       // 商品售价
	Description *string    `json:"description"` // 商品简介
	Status      *int8      `json:"status"`      // 状态：0-下架，1-上架
	CreateTime  *time.Time `json:"createTime"`  // 创建时间
}

// 商品库存信息（简化版）
type ProductStockMO struct {
	Id         *int64     `json:"id"`         // 库存ID
	ProductId  *int64     `json:"productId"`  // 商品ID
	StockNum   *int32     `json:"stockNum"`   // 库存数量
	CreateTime *time.Time `json:"createTime"` // 创建时间
}
