package repository

import (
	"gorm.io/gorm"
	"shopping/user_service/model"
	"shopping/user_service/po"
)

type UserRepository interface {
	// 查找用户登录信息
	GetUserLogin(mo *model.UserMO) (*model.UserMO, error)
	// 查找用户详细信息
	GetUserInfo(mo *model.UserMO) (*model.UserMO, error)
	// 注册新用户登录信息
	CreateUserLogin(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error)
	// 注册新用户详细信息
	CreateUserInfo(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error)
	// 更新用户详细信息
	UpdateUserInfo(db *gorm.DB, mo *model.UserMO) error
	// 更新用户登录信息
	UpdateUserLogin(db *gorm.DB, mo *model.UserMO) error
}

type UserRepositoryImpl struct {
	userLoginPO *po.UserLoginPO
	userInfoPO  *po.UserInfoPO
}

func NewUserRepository(userLoginPO *po.UserLoginPO, userInfoPO *po.UserInfoPO) UserRepository {
	return &UserRepositoryImpl{
		userLoginPO: userLoginPO,
		userInfoPO:  userInfoPO,
	}
}

func (r *UserRepositoryImpl) GetUserLogin(mo *model.UserMO) (*model.UserMO, error) {
	return r.userLoginPO.FindOne(mo)
}

func (r *UserRepositoryImpl) GetUserInfo(mo *model.UserMO) (*model.UserMO, error) {
	return r.userInfoPO.FindOne(mo)
}

func (r *UserRepositoryImpl) CreateUserLogin(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error) {
	return r.userLoginPO.Create(db, mo)

}

func (r *UserRepositoryImpl) CreateUserInfo(db *gorm.DB, mo *model.UserMO) (*model.UserMO, error) {
	return r.userInfoPO.Create(db, mo)
}

func (r *UserRepositoryImpl) UpdateUserInfo(db *gorm.DB, mo *model.UserMO) error {
	return r.userInfoPO.Update(db, mo)
}

func (r *UserRepositoryImpl) UpdateUserLogin(db *gorm.DB, mo *model.UserMO) error {
	return r.userLoginPO.Update(db, mo)
}
