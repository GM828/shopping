package model

import "time"

type UserMO struct {
	UserLoginId *int64     `json:"userLoginId"` // 用户登录ID
	UserName    *string    `json:"userName"`    // 用户名
	Password    *string    `json:"passWord"`    // 密码
	UserInfoId  *int64     `json:"userInfoId"`  // 用户信息ID
	RealName    *string    `json:"realName"`    // 真实姓名
	Phone       *string    `json:"phone"`       // 手机号
	Email       *string    `json:"email"`       // 邮箱
	Gender      *int8      `json:"gender"`      // 性别（0：未知，1：男，2：女）
	Birthday    *time.Time `json:"birthday"`    // 出生日期（日期类型）
	CreateTime  *time.Time `json:"createTime"`  // 创建时间
	UpdateTime  *time.Time `json:"updateTime"`  // 更新时间
}
