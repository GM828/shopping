package po

import (
	"errors"
	"gorm.io/gorm"
	"shopping/database"
	"shopping/product_service/model"
	"time"
)

var DefaultProductPO ProductPO

type ProductPO struct {
	Id          *int64     `json:"id" gorm:"primary_key;column:id"`
	ProductName *string    `json:"productName" gorm:"column:product_name"`
	Category    *string    `json:"category" gorm:"column:category"`
	Price       *float64   `json:"price" gorm:"column:price"`
	Description *string    `json:"description" gorm:"column:description"`
	Status      *int8      `json:"status" gorm:"column:status"`
	CreateTime  *time.Time `json:"createTime" gorm:"column:create_time"`
}

func NewProductPO() *ProductPO {
	return &ProductPO{}
}

func (p *ProductPO) TableName() string {
	return "product"
}

func (p *ProductPO) ToModel() *model.ProductMO {
	return &model.ProductMO{
		Id:          p.Id,
		ProductName: p.ProductName,
		Category:    p.Category,
		Price:       p.Price,
		Description: p.Description,
		Status:      p.Status,
		CreateTime:  p.CreateTime,
	}
}

func ToProductModelList(poList []*ProductPO) []*model.ProductMO {
	list := make([]*model.ProductMO, 0, len(poList))
	for _, po := range poList {
		list = append(list, po.ToModel())
	}
	return list
}

func ToProductPO(mo *model.ProductMO) *ProductPO {
	return &ProductPO{
		Id:          mo.Id,
		ProductName: mo.ProductName,
		Category:    mo.Category,
		Price:       mo.Price,
		Description: mo.Description,
		Status:      mo.Status,
		CreateTime:  mo.CreateTime,
	}
}

func ToProductPOList(modelList []*model.ProductMO) []*ProductPO {
	list := make([]*ProductPO, 0, len(modelList))
	for _, m := range modelList {
		list = append(list, ToProductPO(m))
	}
	return list
}

// BuildDynamicQueryCondition \- 构建动态查询条件
func (p *ProductPO) BuildDynamicQueryCondition(query *gorm.DB, mo *model.ProductMO) *gorm.DB {
	if mo.Id != nil {
		query = query.Where("id = ?", *mo.Id)
	}
	if mo.ProductName != nil {
		query = query.Where("product_name = ?", *mo.ProductName)
	}
	if mo.Category != nil {
		query = query.Where("category = ?", *mo.Category)
	}
	if mo.Status != nil {
		query = query.Where("status = ?", *mo.Status)
	}
	return query
}

// BuildDynamicUpdateTaskMO \- 构建动态更新数据
func (p *ProductPO) BuildDynamicUpdateTaskMO(mo *model.ProductMO) map[string]interface{} {
	data := make(map[string]interface{})

	if mo.ProductName != nil {
		data["product_name"] = *mo.ProductName
	}
	if mo.Category != nil {
		data["category"] = *mo.Category
	}
	if mo.Price != nil {
		data["price"] = *mo.Price
	}
	if mo.Description != nil {
		data["description"] = *mo.Description
	}
	if mo.Status != nil {
		data["status"] = *mo.Status
	}
	// 一般不更新 create_time，这里不放入

	return data
}

func (p *ProductPO) FindOne(mo *model.ProductMO) (*model.ProductMO, error) {
	var po ProductPO
	query := database.DB().Model(&ProductPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *ProductPO) FindList(mo *model.ProductMO) ([]*model.ProductMO, error) {
	var poList []*ProductPO
	query := database.DB().Model(&ProductPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.Find(&poList).Error; err != nil {
		return nil, err
	}
	return ToProductModelList(poList), nil
}

func (p *ProductPO) Create(db *gorm.DB, mo *model.ProductMO) (*model.ProductMO, error) {
	po := ToProductPO(mo)
	if err := db.Create(po).Error; err != nil {
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *ProductPO) Update(db *gorm.DB, mo *model.ProductMO) error {
	data := p.BuildDynamicUpdateTaskMO(mo)
	if len(data) == 0 {
		return nil
	}
	if mo.Id == nil {
		return errors.New("product id is nil")
	}
	if err := db.Model(&ProductPO{}).Where("id = ?", *mo.Id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
