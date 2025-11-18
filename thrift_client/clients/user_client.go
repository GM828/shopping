package clients

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"shopping/thrift_client/pool"
	"shopping/thrift_client/thrift_manager"
	"shopping/thrift_def/gen-go/user"
)

/*
// 单例模式方案

// 全局变量：保存初始化后的user客户端实例
var (
	userClientInstance *user.UserServiceClient
	instanceOnce       sync.Once // 确保客户端实例只被赋值一次
)

// 1. 客户端构造函数（供管理器注册使用）
func userClientFactory(transport thrift.TTransport) interface{} {
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := user.NewUserServiceClientFactory(transport, protocolFactory)

	// 关键：在工厂函数中，将创建的客户端实例赋值给全局变量
	instanceOnce.Do(func() {
		userClientInstance = client
	})

	return client
}

// 2. 注册user服务到ThriftManager
func RegisterUserService(manager *thrift_manager.ThriftManager, addr string) {
	manager.Register("user", addr, userClientFactory)
}

// 3. 业务层直接获取客户端实例（无需通过管理器遍历）
func GetUserClient() (*user.UserServiceClient, bool) {
	// 检查全局变量是否已初始化
	if userClientInstance == nil {
		return nil, false
	}
	return userClientInstance, true
}*/

// 全局连接池（user服务）
var userThriftPool *pool.ThriftPool

// 1. 客户端构造函数（供连接池使用）
func userClientFactory(transport thrift.TTransport) interface{} {
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	return user.NewUserServiceClientFactory(transport, protocolFactory)
}

// 2. 注册user服务到管理器
func RegisterUserService(manager *thrift_manager.ThriftManager, addr string, maxConn int) {
	// 创建连接池工厂
	factory := pool.NewThriftFactory(addr, userClientFactory)
	// 创建连接池
	p, err := pool.NewThriftPool("user", factory, maxConn)
	if err != nil {
		panic("初始化user连接池失败: " + err.Error())
	}
	log.Println("初始化user连接池成功，地址:", addr)
	userThriftPool = p
	// 注册到管理器
	manager.Register("user", p)
}

// 3. 业务层获取客户端的便捷方法
func BorrowUserClient() (*user.UserServiceClient, *pool.WrappedClient, error) {
	wrapped, err := userThriftPool.Borrow()
	if err != nil {
		return nil, nil, err
	}
	// 类型转换（确保安全）
	client, ok := wrapped.Client.(*user.UserServiceClient)
	if !ok {
		return nil, nil, fmt.Errorf("user客户端类型错误")
	}
	return client, wrapped, nil
}

// 4.业务层归还客户端的便捷方法
func ReturnUserClient(wrapped *pool.WrappedClient) {
	userThriftPool.Return(wrapped)
}
