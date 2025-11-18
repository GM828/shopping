package service

import (
	"gorm.io/gorm"
	"shopping/database"
	"shopping/user_service/model"
	"shopping/user_service/repository"
)

type UserService interface {
	// 查询用户登录信息
	GetUserLogin(mo *model.UserMO) (*model.UserMO, error)
	// 查询用户详细信息
	GetUserInfo(mo *model.UserMO) (*model.UserMO, error)
	// 注册新用户
	CreateUser(mo *model.UserMO) (*model.UserMO, error)
	// 更新用户登录信息
	UpdateUserLogin(mo *model.UserMO) error
	// 更新用户详细信息
	UpdateUserInfo(mo *model.UserMO) error
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (s *UserServiceImpl) GetUserLogin(mo *model.UserMO) (*model.UserMO, error) {
	userMO := model.UserMO{
		UserLoginId: mo.UserLoginId,
		UserName:    mo.UserName,
		Password:    mo.Password,
		Email:       mo.Email,
	}
	return s.userRepository.GetUserLogin(&userMO)
}

func (s *UserServiceImpl) GetUserInfo(mo *model.UserMO) (*model.UserMO, error) {
	userMO := model.UserMO{
		UserLoginId: mo.UserLoginId,
	}
	return s.userRepository.GetUserInfo(&userMO)
}

func (s *UserServiceImpl) CreateUser(mo *model.UserMO) (*model.UserMO, error) {
	registerMO := model.UserMO{
		UserName: mo.UserName,
		Password: mo.Password,
		Email:    mo.Email,
	}
	// 事务处理
	err := database.DB().Transaction(func(tx *gorm.DB) error {
		userLogin, err := s.userRepository.CreateUserLogin(tx, &registerMO)
		if err != nil {
			return err
		}
		registerMO.UserLoginId = userLogin.UserLoginId
		registerMO.RealName = mo.RealName
		registerMO.Phone = mo.Phone
		registerMO.Gender = mo.Gender
		registerMO.Birthday = mo.Birthday
		info, err := s.userRepository.CreateUserInfo(tx, &registerMO)
		if err != nil {
			return err
		}
		registerMO.UserInfoId = info.UserInfoId
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &registerMO, nil
}

func (s *UserServiceImpl) UpdateUserLogin(mo *model.UserMO) error {
	err := s.userRepository.UpdateUserLogin(database.DB(), mo)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) UpdateUserInfo(mo *model.UserMO) error {
	err := s.userRepository.UpdateUserInfo(database.DB(), mo)
	if err != nil {
		return err
	}
	return nil
}
