# 前言

在[之前的文章](https://blog.csdn.net/2201_75669520/article/details/154449136)当中学习了从0开始在Go当中使用Apache Thrift框架，这篇文章主要是记录我在一个大型分布式商城项目当中，使用Thrift作为RPC框架的实战记录，主要记录方向是**高并发优化**这一块

文章地址：



# 1.登录态

在Java代码当中，登录态这个东西，我们大多情况下是自己保存登录态信息，比如保存到redis当中，然后把token交给前端，或者说使用jwt。实际开发的过程当中，登录态这个东西，都是交给网关部门去做的，如果说作为业务开发人员，只需要从请求头当中获取登录态信息即可

例如这样

```go
garageID := ctx.GetHeader("X-Garage-ID")
userID := ctx.GetHeader("X-User-ID")
if garageID == "" || userID == "" {
    // 若缺失，说明网关校验有问题，直接返回401
    ctx.JSON(401, common.Error("未获取到用户信息"))
    return
}
```

我们自己开发项目的时候，模拟一下就好，在接口工具当中，填写header即可

<img width="994" height="288" alt="image" src="https://github.com/user-attachments/assets/0a0e1cf8-75f5-499e-9fb0-118507e8321e" />




# 2.生成gen-go文件当中的字段

```
// 登录请求参数，对应 UserLoginRequestDTO
struct UserLoginRequest {
    1: optional i64 userLoginId,   // 用户登录ID
    2: optional string userName,   // 用户名
    3: required string email,      // 邮箱
    4: required string password    // 密码
}
```

这是我的IDL

```go
type UserLoginRequest struct {
    UserLoginId *int64 `thrift:"userLoginId,1" db:"userLoginId" json:"userLoginId,omitempty"`
    UserName *string `thrift:"userName,2" db:"userName" json:"userName,omitempty"`
    Email string `thrift:"email,3,required" db:"email" json:"email"`
    Password string `thrift:"password,4,required" db:"password" json:"password"`
}
```

这是生成的结构体代码

实际使用过程中，因为Go语言是显示指针，和Java不一样，因此我们大部分接收都是用指针接收的，但是他这里生成的代码不是指针，原因是它们在 Thrift IDL 中被声明为 `required`。在 Go 中，`required` 字段会被生成为非指针类型（如 `string`），而 `optional` 字段会被生成为指针类型（如 `*string`）

这样设计的目的是通过 `optional` 明确字段的可选性，生成指针类型后可以通过 `nil` 表示字段未设置，而 `required` 字段强制要求必须赋值，因此直接使用非指针类型来保证字段的存在性。

解决方案有两种

**1.根据公司项目要求，所有字段统一声明为optional**

如果这样做的话，字段的存在性校验就需要放到具体的Service当中去判断了，我代码当中目前主推这种做法，因为我的各个service当中，都在manager层进行了字段的存在性校验



**2.完善服务端代码：实现IDL当中定义的接口的时候，把字段的值进行取指操作，代码如下**

```go
func (h *UserHandler) Login(ctx context.Context, request *user.UserLoginRequest) (*user.UserResponse, error) {
    loginReq := &dto.UserLoginRequestDTO{
       UserLoginId: request.UserLoginId,
       UserName:    request.UserName,
       Email:       &request.Email,
       Password:    &request.Password,
    }
```

可以看到，我在接收客户端传过来的*user.UserLoginRequest时，需要用服务端的&dto.UserLoginRequestDTO去接受，就类似接口层的Controller一样，然后去调用Manager的方法，此时，email和password这两个设置为required的字段，就进行取指操作就好

后续文章：https://blog.csdn.net/2201_75669520/article/details/155005840
其他文章：https://blog.csdn.net/2201_75669520?type=blog


# 新增Kitex框架实战

新增product_service，使用Kitex框架进行客户端和服务端连接

具体文章：https://mp.csdn.net/mp_blog/creation/editor/155242379
