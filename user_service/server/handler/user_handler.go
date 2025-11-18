package handler

import (
	"context"
	"log"
	"shopping/thrift_def/gen-go/user"
	"shopping/user_service/dto"
	"shopping/user_service/manager"
)

type UserHandler struct {
	UserManager manager.UserManager
}

// Ping 心跳检测
func (s *UserHandler) Ping(ctx context.Context) error {
	// 无需任何业务逻辑，直接返回成功
	// 可选：记录心跳日志（建议只在调试时开启）
	log.Println("[心跳检测] Ping received")
	return nil
}

// 登录
func (h *UserHandler) Login(ctx context.Context, request *user.UserLoginRequest) (*user.UserResponse, error) {
	loginReq := &dto.UserLoginRequestDTO{
		UserLoginId: request.UserLoginId,
		UserName:    request.UserName,
		Email:       request.Email,
		Password:    request.Password,
	}
	respDTO, err := h.UserManager.UserLogin(loginReq)
	if err != nil {
		code := int32(500)
		return nil, &user.UserException{Message: err.Error(), Code: &code}
	}
	return &user.UserResponse{
		UserName: respDTO.UserName,
		Phone:    respDTO.Phone,
		Email:    respDTO.Email,
		Gender:   respDTO.Gender,
		Birthday: respDTO.Birthday,
	}, nil
}

// 注册
func (h *UserHandler) Register(ctx context.Context, request *user.UserRegisterRequest) (*user.CommonResponse, error) {
	regReq := &dto.UserRegisterRequestDTO{
		UserName: request.UserName,
		Email:    request.Email,
		Password: request.Password,
		RealName: request.RealName,
		Phone:    request.Phone,
		Gender:   request.Gender,
		Birthday: request.Birthday,
	}
	err := h.UserManager.UserRegister(regReq)
	if err != nil {
		code := int32(500)
		return nil, &user.UserException{Message: "注册失败：" + err.Error(), Code: &code}
	}
	return &user.CommonResponse{Message: "注册成功"}, nil
}

// 修改密码
func (h *UserHandler) UpdatePassword(ctx context.Context, request *user.UserUpdateLoginRequest) (*user.CommonResponse, error) {
	updateReq := &dto.UserUpdateLoginRequestDTO{
		UserName:    request.UserName,
		Email:       request.Email,
		Password:    request.Password,
		NewPassword: request.NewPassword_,
	}
	err := h.UserManager.UpdatePassword(updateReq)
	if err != nil {
		code := int32(500)
		return nil, &user.UserException{Message: "修改密码失败：" + err.Error(), Code: &code}
	}
	return &user.CommonResponse{Message: "修改密码成功"}, nil
}

// 修改用户信息
func (h *UserHandler) UpdateUserInfo(ctx context.Context, request *user.UserUpdateInfoRequest) (*user.CommonResponse, error) {
	updateReq := &dto.UserUpdateInfoRequestDTO{
		UserInfoId: request.UserInfoId,
		RealName:   request.RealName,
		Phone:      request.Phone,
		Gender:     request.Gender,
		Birthday:   request.Birthday,
	}
	err := h.UserManager.UpdateUserInfo(updateReq)
	if err != nil {
		code := int32(500)
		return nil, &user.UserException{Message: "修改用户信息失败：" + err.Error(), Code: &code}
	}
	return &user.CommonResponse{Message: "修改用户信息成功"}, nil
}
