package po

import (
	"errors"
	"gorm.io/gorm"
	"shopping/database"
	"shopping/user_service/model"
)

var DefaultUserPO UserLoginPO

type UserLoginPO struct {
	Id       *int64  `json:"userLoginId" gorm:"primary_key;column:id"`
	UserName *string `json:"userName" gorm:"column:username"`
	Password *string `json:"passWord" gorm:"column:password"`
	Email    *string `json:"email" gorm:"column:email"` // 邮箱
}

func NewUserLoginPO() *UserLoginPO {
	return &UserLoginPO{}
}

// TableName 返回数据库表名
func (p *UserLoginPO) TableName() string {
	return "user_login"
}

func (p *UserLoginPO) ToModel() *model.UserMO {
	return &model.UserMO{
		UserLoginId: p.Id,
		UserName:    p.UserName,
		Password:    p.Password,
	}
}

func ToModelList(poList []*UserLoginPO) []*model.UserMO {
	modelList := make([]*model.UserMO, 0)
	for _, po := range poList {
		modelList = append(modelList, po.ToModel())
	}
	return modelList
}

func ToUserLoginPO(model *model.UserMO) *UserLoginPO {
	return &UserLoginPO{
		Id:       model.UserLoginId,
		UserName: model.UserName,
		Password: model.Password,
		Email:    model.Email,
	}
}

func ToUserLoginPOList(modelList []*model.UserMO) []*UserLoginPO {
	poList := make([]*UserLoginPO, 0)
	for _, model := range modelList {
		poList = append(poList, ToUserLoginPO(model))
	}
	return poList
}

// BuildDynamicQueryCondition 构建动态查询条件
func (p *UserLoginPO) BuildDynamicQueryCondition(query *gorm.DB, mo *model.UserMO) *gorm.DB {
	if mo.UserLoginId != nil {
		query = query.Where("id = ?", *mo.UserLoginId)
	}
	if mo.UserName != nil {
		query = query.Where("username = ?", *mo.UserName)
	}
	if mo.Email != nil {
		query = query.Where("email = ?", *mo.Email)
	}
	return query
}

// BuildDynamicUpdateTaskMO 构建动态更新数据
func (p *UserLoginPO) BuildDynamicUpdateTaskMO(mo *model.UserMO) map[string]interface{} {
	updateData := make(map[string]interface{})
	if mo.UserName != nil {
		updateData["username"] = *mo.UserName
	}
	if mo.Password != nil {
		updateData["password"] = *mo.Password
	}
	return updateData
}

func (p *UserLoginPO) FindOne(mo *model.UserMO) (*model.UserMO, error) {
	var po UserLoginPO
	query := database.DB().Model(&UserLoginPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *UserLoginPO) FindList(mo *model.UserMO) ([]*model.UserMO, error) {
	var poList []*UserLoginPO
	query := database.DB().Model(&UserLoginPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.Find(&poList).Error; err != nil {
		return nil, err
	}
	return ToModelList(poList), nil
}

func (p *UserLoginPO) Create(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error) {
	po := ToUserLoginPO(mo)
	if err := db.Create(po).Error; err != nil {
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *UserLoginPO) Update(db *gorm.DB, mo *model.UserMO) error {
	updateData := p.BuildDynamicUpdateTaskMO(mo)
	if err := db.Model(&UserLoginPO{}).Where("id = ?", *mo.UserLoginId).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}
