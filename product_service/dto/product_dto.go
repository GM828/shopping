package dto

// 商品核心信息
type ProductRequestDTO struct {
	Id          *int64   `json:"id"`          // 商品ID
	ProductName *string  `json:"productName"` // 商品名称
	Category    *string  `json:"category"`    // 商品分类（直接存储名称）
	Price       *float64 `json:"price"`       // 商品售价
	Description *string  `json:"description"` // 商品简介
	Status      *int8    `json:"status"`      // 状态：0-下架，1-上架
	CreateTime  *string  `json:"createTime"`  // 创建时间
}

type ProductResponseDTO struct {
	Id          *int64   `json:"id"`          // 商品ID
	ProductName *string  `json:"productName"` // 商品名称
	Category    *string  `json:"category"`    // 商品分类（直接存储名称）
	Price       *float64 `json:"price"`       // 商品售价
	Description *string  `json:"description"` // 商品简介
	Status      *int8    `json:"status"`      // 状态：0-下架，1-上架
	CreateTime  *string  `json:"createTime"`  // 创建时间
	StockNum    *int32   `json:"stockNum"`    // 库存数量
}
