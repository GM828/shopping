package dto

import (
	"shopping/user_service/model"
	"shopping/util"
)

type UserLoginRequestDTO struct {
	UserLoginId *int64  `json:"userLoginId"`
	UserName    *string `json:"userName"`
	Email       *string `json:"email"`
	Password    *string `json:"passWord"`
}

type UserRegisterRequestDTO struct {
	UserName *string `json:"userName"`
	Email    *string `json:"email"`
	Password *string `json:"passWord"`
	RealName *string `json:"realName"` // 真实姓名
	Phone    *string `json:"phone"`    // 手机号
	Gender   *int8   `json:"gender"`   // 性别（0：未知，1：男，2：女）
	Birthday *string `json:"birthday"` // 出生日期（日期类型）
}

type UserUpdateLoginRequestDTO struct {
	UserName    *string `json:"userName"`
	Email       *string `json:"email"`
	Password    *string `json:"passWord"`
	NewPassword *string `json:"newPassWord"`
}

type UserUpdateInfoRequestDTO struct {
	UserInfoId *int64  `json:"userInfoId"` // 用户信息ID
	RealName   *string `json:"realName"`   // 真实姓名
	Phone      *string `json:"phone"`      // 手机号
	Gender     *int8   `json:"gender"`     // 性别（0：未知，1：男，2：女）
	Birthday   *string `json:"birthday"`   // 出生日期（日期类型）
}

type UserResponseDTO struct {
	UserName *string `json:"userName"` // 用户名
	Phone    *string `json:"phone"`    // 手机号
	Email    *string `json:"email"`    // 邮箱
	Gender   *int8   `json:"gender"`   // 性别（0：未知，1：男，2：女）
	Birthday *string `json:"birthday"` // 出生日期（日期类型）
}

func UserRequestDTOToMo(userRequestDTO *UserLoginRequestDTO) *model.UserMO {
	return &model.UserMO{
		UserLoginId: userRequestDTO.UserLoginId,
		UserName:    userRequestDTO.UserName,
		Password:    userRequestDTO.Password,
		Email:       userRequestDTO.Email,
	}
}

func UserMoToResponseDTO(userMo *model.UserMO) *UserResponseDTO {
	birthday := util.DateUtil.FormatDateByCustomLayout(userMo.Birthday, util.DateLayout.YYYY_MM_DD_HH_MM_SS)
	dto := &UserResponseDTO{
		UserName: userMo.UserName,
		Phone:    userMo.Phone,
		Email:    userMo.Email,
		Gender:   userMo.Gender,
		Birthday: &birthday,
	}
	return dto
}
