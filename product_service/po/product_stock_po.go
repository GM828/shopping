package po

import (
	"errors"
	"gorm.io/gorm"
	"shopping/database"
	"shopping/product_service/model"
	"time"
)

var DefaultProductStockPO ProductStockPO

type ProductStockPO struct {
	Id         *int64     `json:"id" gorm:"primary_key;column:id"`
	ProductId  *int64     `json:"productId" gorm:"column:product_id"`
	StockNum   *int32     `json:"stockNum" gorm:"column:stock_num"`
	CreateTime *time.Time `json:"createTime" gorm:"column:create_time"`
}

func NewProductStockPO() *ProductStockPO {
	return &ProductStockPO{}
}

func (p *ProductStockPO) TableName() string {
	return "product_stock"
}

func (p *ProductStockPO) ToModel() *model.ProductStockMO {
	return &model.ProductStockMO{
		Id:         p.Id,
		ProductId:  p.ProductId,
		StockNum:   p.StockNum,
		CreateTime: p.CreateTime,
	}
}

func ToProductStockModelList(poList []*ProductStockPO) []*model.ProductStockMO {
	list := make([]*model.ProductStockMO, 0, len(poList))
	for _, po := range poList {
		list = append(list, po.ToModel())
	}
	return list
}

func ToProductStockPO(mo *model.ProductStockMO) *ProductStockPO {
	return &ProductStockPO{
		Id:         mo.Id,
		ProductId:  mo.ProductId,
		StockNum:   mo.StockNum,
		CreateTime: mo.CreateTime,
	}
}

func ToProductStockPOList(modelList []*model.ProductStockMO) []*ProductStockPO {
	list := make([]*ProductStockPO, 0, len(modelList))
	for _, m := range modelList {
		list = append(list, ToProductStockPO(m))
	}
	return list
}

// BuildDynamicQueryCondition \- 构建动态查询条件
func (p *ProductStockPO) BuildDynamicQueryCondition(query *gorm.DB, mo *model.ProductStockMO) *gorm.DB {
	if mo.Id != nil {
		query = query.Where("id = ?", *mo.Id)
	}
	if mo.ProductId != nil {
		query = query.Where("product_id = ?", *mo.ProductId)
	}
	return query
}

// BuildDynamicUpdateTaskMO \- 构建动态更新数据
func (p *ProductStockPO) BuildDynamicUpdateTaskMO(mo *model.ProductStockMO) map[string]interface{} {
	data := make(map[string]interface{})
	if mo.StockNum != nil {
		data["stock_num"] = *mo.StockNum
	}
	// 一般不更新 create_time、product_id
	return data
}

func (p *ProductStockPO) FindOne(mo *model.ProductStockMO) (*model.ProductStockMO, error) {
	var po ProductStockPO
	query := database.DB().Model(&ProductStockPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *ProductStockPO) FindList(mo *model.ProductStockMO) ([]*model.ProductStockMO, error) {
	var poList []*ProductStockPO
	query := database.DB().Model(&ProductStockPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.Find(&poList).Error; err != nil {
		return nil, err
	}
	return ToProductStockModelList(poList), nil
}

func (p *ProductStockPO) Create(db *gorm.DB, mo *model.ProductStockMO) (*model.ProductStockMO, error) {
	po := ToProductStockPO(mo)
	if err := db.Create(po).Error; err != nil {
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *ProductStockPO) Update(db *gorm.DB, mo *model.ProductStockMO) error {
	data := p.BuildDynamicUpdateTaskMO(mo)
	if len(data) == 0 {
		return nil
	}
	if mo.Id == nil {
		return errors.New("product stock id is nil")
	}
	if err := db.Model(&ProductStockPO{}).Where("id = ?", *mo.Id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}
