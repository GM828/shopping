namespace go user

// 登录请求参数，对应 UserLoginRequestDTO
struct UserLoginRequest {
    1: optional i64 userLoginId,   // 用户登录ID
    2: optional string userName,   // 用户名
    3: optional string email,      // 邮箱
    4: optional string password    // 密码
}

// 注册请求参数，对应 UserRegisterRequestDTO
struct UserRegisterRequest {
    1: optional string userName,   // 用户名
    2: optional string email,      // 邮箱
    3: optional string password,   // 密码
    4: optional string realName,   // 真实姓名
    5: optional string phone,      // 手机号
    6: optional i8 gender,         // 性别（0：未知，1：男，2：女）
    7: optional string birthday    // 出生日期（如 "2000-01-01"）
}

// 修改密码请求参数，对应 UserUpdateLoginRequestDTO
struct UserUpdateLoginRequest {
    1: optional string userName,   // 用户名
    2: optional string email,      // 邮箱
    3: optional string password,   // 旧密码
    4: optional string newPassword // 新密码
}

// 修改用户信息请求参数，对应 UserUpdateInfoRequestDTO
struct UserUpdateInfoRequest {
    1: optional i64 userInfoId,    // 用户信息ID
    2: optional string realName,   // 真实姓名
    3: optional string phone,      // 手机号
    4: optional i8 gender,         // 性别（0：未知，1：男，2：女）
    5: optional string birthday    // 出生日期（如 "2000-01-01"）
}

// 通用响应
struct CommonResponse {
    1: required string message     // 响应消息
}

// 用户响应信息
struct UserResponse {
    1: optional string userName,   // 用户名
    2: optional string phone,      // 手机号
    3: optional string email,      // 邮箱
    4: optional i8 gender,         // 性别（0：未知，1：男，2：女）
    5: optional string birthday    // 出生日期（如 "2000-01-01"）
}

// 业务异常
exception UserException {
    1: required string message,    // 错误详情
    2: optional i32 code           // 错误码
}

// 用户服务接口
service UserService {
    UserResponse login(1: UserLoginRequest request) throws (1: UserException ex),                // 登录
    CommonResponse register(1: UserRegisterRequest request) throws (1: UserException ex),        // 注册
    CommonResponse updatePassword(1: UserUpdateLoginRequest request) throws (1: UserException ex), // 修改密码
    CommonResponse updateUserInfo(1: UserUpdateInfoRequest request) throws (1: UserException ex) // 修改用户信息

    // 心跳检测方法
    void Ping()
}
