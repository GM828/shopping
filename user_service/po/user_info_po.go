package po

import (
	"errors"
	"gorm.io/gorm"
	"shopping/database"
	"shopping/user_service/model"
	"shopping/util"
	"time"
)

type UserInfoPO struct {
	Id          *int64     `json:"id" gorm:"primary_key;column:id"`           // 主键ID
	UserLoginId *int64     `json:"userId" gorm:"column:user_login_id"`        // 关联用户登录的ID（外键）
	RealName    *string    `json:"realName" gorm:"column:real_name"`          // 真实姓名
	Phone       *string    `json:"phone" gorm:"column:phone"`                 // 手机号
	Gender      *int8      `json:"gender" gorm:"column:gender"`               // 性别（0：未知，1：男，2：女）
	Birthday    *time.Time `json:"birthday" gorm:"column:birthday;type:date"` // 出生日期（日期类型）
	CreateTime  *time.Time `json:"createTime" gorm:"column:create_time"`      // 创建时间
	UpdateTime  *time.Time `json:"updateTime" gorm:"column:update_time"`      // 更新时间
}

func NewUserInfoPO() *UserInfoPO {
	return &UserInfoPO{}
}

// TableName 指定数据库表名
func (p *UserInfoPO) TableName() string {
	return "user_info"
}

func (p *UserInfoPO) ToModel() *model.UserMO {
	return &model.UserMO{
		UserInfoId:  p.Id,
		UserLoginId: p.UserLoginId,
		RealName:    p.RealName,
		Phone:       p.Phone,
		Gender:      p.Gender,
		Birthday:    p.Birthday,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
	}
}

func ToUserInfoPO(mo *model.UserMO) *UserInfoPO {
	return &UserInfoPO{
		Id:          mo.UserInfoId,
		UserLoginId: mo.UserLoginId,
		RealName:    mo.RealName,
		Phone:       mo.Phone,
		Gender:      mo.Gender,
		Birthday:    mo.Birthday,
	}
}

// BuildDynamicQueryCondition 构建动态查询条件
func (p *UserInfoPO) BuildDynamicQueryCondition(query *gorm.DB, mo *model.UserMO) *gorm.DB {
	if mo.UserInfoId != nil {
		query = query.Where("id = ?", *mo.UserInfoId)
	}
	if mo.UserLoginId != nil {
		query = query.Where("user_login_id = ?", *mo.UserLoginId)
	}
	return query
}

// BuildDynamicUpdateTaskMO 构建动态更新数据
func (p *UserInfoPO) BuildDynamicUpdateTaskMO(mo *model.UserMO) map[string]interface{} {
	updateData := make(map[string]interface{})
	if mo.RealName != nil {
		updateData["real_name"] = *mo.RealName
	}
	if mo.Phone != nil {
		updateData["phone"] = *mo.Phone
	}
	if mo.Gender != nil {
		updateData["gender"] = *mo.Gender
	}
	if mo.Birthday != nil {
		updateData["birthday"] = *mo.Birthday
	}
	return updateData
}

func (p *UserInfoPO) FindOne(mo *model.UserMO) (*model.UserMO, error) {
	var po UserInfoPO
	query := database.DB().Model(&UserInfoPO{})
	query = p.BuildDynamicQueryCondition(query, mo)
	if err := query.First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *UserInfoPO) Create(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error) {
	po := ToUserInfoPO(mo)
	po.CreateTime = util.DateUtil.TimePtr(time.Now())
	po.UpdateTime = util.DateUtil.TimePtr(time.Now())
	if err := db.Create(po).Error; err != nil {
		return nil, err
	}
	return po.ToModel(), nil
}

func (p *UserInfoPO) Update(db *gorm.DB, mo *model.UserMO) error {
	updateData := p.BuildDynamicUpdateTaskMO(mo)
	updateData["update_time"] = time.Now()
	if err := db.Model(&UserInfoPO{}).Where("id = ?", *mo.UserInfoId).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}
