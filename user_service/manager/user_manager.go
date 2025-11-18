package manager

import (
	"errors"
	"shopping/user_service/dto"
	"shopping/user_service/model"
	"shopping/user_service/service"
	"shopping/util"
)

type UserManager interface {
	// 用户登录
	UserLogin(requestDTO *dto.UserLoginRequestDTO) (*dto.UserResponseDTO, error)
	// 用户注册
	UserRegister(requestDTO *dto.UserRegisterRequestDTO) error
	// 更新用户密码
	UpdatePassword(requestDTO *dto.UserUpdateLoginRequestDTO) error
	// 更新用户信息
	UpdateUserInfo(requestDTO *dto.UserUpdateInfoRequestDTO) error
}
type UserManagerImpl struct {
	userService service.UserService
}

func NewUserManager(userService service.UserService) UserManager {
	return &UserManagerImpl{
		userService: userService,
	}
}

func (m *UserManagerImpl) UserLogin(requestDTO *dto.UserLoginRequestDTO) (*dto.UserResponseDTO, error) {
	// 参数校验
	if requestDTO.Email == nil {
		return nil, errors.New("邮箱不能为空")
	}
	if requestDTO.Password == nil {
		return nil, errors.New("密码不能为空")
	}

	// 查询用户是否存在
	userLogin, err := m.userService.GetUserLogin(dto.UserRequestDTOToMo(&dto.UserLoginRequestDTO{
		UserLoginId: requestDTO.UserLoginId,
		Email:       requestDTO.Email,
	}))
	if err != nil {
		return nil, err
	}
	if userLogin == nil {
		return nil, errors.New("用户不存在")
	}

	// 校验密码
	if *userLogin.Password != *requestDTO.Password {
		return nil, errors.New("密码错误")
	}

	// 查询用户的额外信息
	userInfo, err := m.userService.GetUserInfo(&model.UserMO{
		UserLoginId: userLogin.UserLoginId,
	})
	if err != nil {
		return nil, err
	}

	// 组装用户信息
	userInfo.UserName = userLogin.UserName
	userInfo.Email = userLogin.Email

	// TODO: 生成token...
	// token := m.tokenService.GenerateToken(userLogin.UserLoginId)

	return dto.UserMoToResponseDTO(userInfo), nil
}

func (m *UserManagerImpl) UserRegister(requestDTO *dto.UserRegisterRequestDTO) error {
	// 参数校验
	if requestDTO.Email == nil {
		return errors.New("邮箱不能为空")
	}
	if requestDTO.Password == nil {
		return errors.New("密码不能为空")
	}
	if requestDTO.UserName == nil {
		return errors.New("用户名不能为空")
	}

	// 查询用户是否存在
	existingUser, err := m.userService.GetUserLogin(&model.UserMO{
		Email: requestDTO.Email,
	})
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("当前邮箱已被注册，请重新登录或修改密码")
	}

	// 创建用户登录信息
	registerMO := &model.UserMO{
		UserName: requestDTO.UserName,
		Email:    requestDTO.Email,
		Password: requestDTO.Password,
		RealName: requestDTO.RealName,
		Phone:    requestDTO.Phone,
		Gender:   requestDTO.Gender,
	}
	birthday, err := util.DateUtil.ParseStandard(requestDTO.Birthday)
	if err != nil {
		return errors.New("生日格式错误")
	}
	registerMO.Birthday = birthday

	// 调用服务层创建用户
	_, err = m.userService.CreateUser(registerMO)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserManagerImpl) UpdatePassword(requestDTO *dto.UserUpdateLoginRequestDTO) error {
	// 参数校验
	if requestDTO.Email == nil {
		return errors.New("邮箱不能为空")
	}
	if requestDTO.Password == nil {
		return errors.New("密码不能为空")
	}
	if requestDTO.NewPassword == nil {
		return errors.New("新密码不能为空")
	}

	// 查询用户是否存在
	userLogin, err := m.userService.GetUserLogin(dto.UserRequestDTOToMo(&dto.UserLoginRequestDTO{
		Email: requestDTO.Email,
	}))
	if err != nil {
		return err
	}
	if userLogin == nil {
		return errors.New("当前用户不存在，无法修改密码")
	}

	// 对比旧密码
	if *userLogin.Password != *requestDTO.Password {
		return errors.New("旧密码错误，无法修改密码")
	}

	// 更新密码
	err = m.userService.UpdateUserLogin(&model.UserMO{
		UserLoginId: userLogin.UserLoginId,
		Password:    requestDTO.NewPassword,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *UserManagerImpl) UpdateUserInfo(requestDTO *dto.UserUpdateInfoRequestDTO) error {
	// 参数校验
	if requestDTO.UserInfoId == nil {
		return errors.New("用户信息ID不能为空")
	}

	// 更新用户信息
	birthday, err := util.DateUtil.ParseStandard(requestDTO.Birthday)
	if err != nil {
		return errors.New("生日格式错误")
	}
	err = m.userService.UpdateUserInfo(&model.UserMO{
		UserInfoId: requestDTO.UserInfoId,
		RealName:   requestDTO.RealName,
		Phone:      requestDTO.Phone,
		Gender:     requestDTO.Gender,
		Birthday:   birthday,
	})
	if err != nil {
		return err
	}
	return nil
}
