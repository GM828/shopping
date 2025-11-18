package pool

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jolestar/go-commons-pool/v2"
	"log"
	"shopping/thrift_def/gen-go/user"
	"time"
)

// ------------------------------
// 连接池对象包装类（显式保存transport）
// ------------------------------
type WrappedClient struct {
	Client    interface{}       // 实际Thrift客户端（如*user.UserServiceClient）
	Transport thrift.TTransport // 对应的transport，用于关闭连接
}

// ------------------------------
// 连接池工厂（实现PooledObjectFactory接口）
// ------------------------------
type ThriftFactory struct {
	addr    string                              // 服务地址（如"localhost:9090"）
	factory func(thrift.TTransport) interface{} // 客户端构造函数
}

func NewThriftFactory(addr string, factory func(thrift.TTransport) interface{}) *ThriftFactory {
	return &ThriftFactory{
		addr:    addr,
		factory: factory,
	}
}

// MakeObject 创建连接（池化对象）
func (f *ThriftFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	// 1. 创建并打开transport
	socket, err := thrift.NewTSocket(f.addr)
	if err != nil {
		return nil, fmt.Errorf("创建transport失败: %w", err)
	}
	if err := socket.Open(); err != nil {
		return nil, fmt.Errorf("连接服务失败: %w", err)
	}

	// 2. 优化：长连接包装为 TFramedTransport（关键变更）
	transport := thrift.NewTFramedTransport(socket)

	// 3. 创建客户端并包装（显式保存transport）
	client := f.factory(transport)
	// 关键：创建指针类型的包装对象
	wrapped := &WrappedClient{
		Client:    client,
		Transport: transport,
	}

	return pool.NewPooledObject(wrapped), nil
}

// DestroyObject 销毁连接（关闭transport）
func (f *ThriftFactory) DestroyObject(ctx context.Context, obj *pool.PooledObject) error {
	wrapped, ok := obj.Object.(*WrappedClient)
	if !ok {
		return fmt.Errorf("销毁失败：对象类型错误（%T）", obj.Object)
	}
	// 关闭transport释放连接
	if wrapped.Transport.IsOpen() {
		return wrapped.Transport.Close()
	}
	return nil
}

// ValidateObject 验证连接是否有效（可选）
func (f *ThriftFactory) ValidateObject(ctx context.Context, obj *pool.PooledObject) bool {
	wrapped, ok := obj.Object.(*WrappedClient)
	if !ok {
		return false
	}

	// 1. 基础检查：连接是否打开
	if !wrapped.Transport.IsOpen() {
		log.Printf("[连接池] ValidateObject 失败: 连接未打开\n")
		return false
	}

	// 2. 心跳检查（核心）
	// 注意：需要服务端支持 Ping 方法
	if client, ok := wrapped.Client.(*user.UserServiceClient); ok {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := client.Ping(ctx); err != nil {
			log.Printf("[连接池] ValidateObject 失败: 心跳检测失败: %v\n", err)
			return false
		}
	}

	// 连接已打开则视为有效
	return true
}

// ActivateObject 激活对象（从池取出时调用，无需操作）
func (f *ThriftFactory) ActivateObject(ctx context.Context, obj *pool.PooledObject) error {
	return nil
}

// PassivateObject 钝化对象（放回池时调用，无需操作）
func (f *ThriftFactory) PassivateObject(ctx context.Context, obj *pool.PooledObject) error {
	return nil
}

// ------------------------------
// 连接池封装（对外提供 borrow/return 方法）
// ------------------------------
type ThriftPool struct {
	pool     *pool.ObjectPool
	poolName string
}

func NewThriftPool(poolName string, factory *ThriftFactory, maxTotal int) (*ThriftPool, error) {
	config := &pool.ObjectPoolConfig{
		MaxTotal:                maxTotal,         // 最大连接数
		MaxIdle:                 maxTotal,         // 最大空闲连接
		MinIdle:                 1,                // 最小空闲连接
		TimeBetweenEvictionRuns: 30 * time.Second, // 闲置连接检测周期
		MinEvictableIdleTime:    5 * time.Minute,  // 连接最大闲置时间（超过则回收）

		// 优化：启用连接验证
		TestOnBorrow:  true, // ← Borrow 时自动验证
		TestOnReturn:  true, // ← Return 时自动验证
		TestWhileIdle: true, // ← 闲置时定期验证
	}
	objPool := pool.NewObjectPool(context.Background(), factory, config)
	return &ThriftPool{
		pool:     objPool,
		poolName: poolName,
	}, nil
}

// Borrow 获取客户端（业务层调用）
func (p *ThriftPool) Borrow() (*WrappedClient, error) {
	obj, err := p.pool.BorrowObject(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取连接失败: %w", err)
	}
	return obj.(*WrappedClient), nil
}

// Return 归还客户端（业务层调用）
func (p *ThriftPool) Return(client *WrappedClient) {
	// 长连接方案：归还，放回连接池
	_ = p.pool.ReturnObject(context.Background(), client)
	// 短连接方案：不归还,直接销毁连接(每次都用新连接)
	// _ = p.pool.InvalidateObject(context.Background(), client)
}

// Close 关闭连接池（管理器调用）
func (p *ThriftPool) Close() {
	p.pool.Close(context.Background())
}
