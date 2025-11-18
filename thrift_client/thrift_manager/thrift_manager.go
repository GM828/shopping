package thrift_manager

import (
	"shopping/thrift_client/pool"
	"sync"
)

// ThriftManager 管理所有服务的连接池
type ThriftManager struct {
	pools []*pool.ThriftPool
	mu    sync.Mutex
}

func NewThriftManager() *ThriftManager {
	return &ThriftManager{
		pools: make([]*pool.ThriftPool, 0),
	}
}

// Register 注册连接池
func (m *ThriftManager) Register(serviceName string, p *pool.ThriftPool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pools = append(m.pools, p)
}

// Init 初始化（连接池创建时已初始化，此处仅日志提示）
func (m *ThriftManager) Init() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 可添加初始化日志
	return nil
}

// Close 批量关闭所有连接池
func (m *ThriftManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var lastErr error
	for _, p := range m.pools {
		p.Close()
	}
	return lastErr
}

/*
单例模式方案

// Transport 封装单个服务的transport和客户端（客户端类型由具体服务决定，这里用interface{}兼容）
type Transport struct {
	ServiceName string            // 服务名（如"user"、"product"）
	Addr        string            // 服务地址（如"localhost:9090"）
	Transport   thrift.TTransport // 传输层
	Client      interface{}       // 具体服务的客户端实例（如*user.UserServiceClient）
}

// ThriftManager 管理所有Thrift服务的连接
type ThriftManager struct {
	transports []*Transport // 注册的所有服务连接
	mu         sync.Mutex   // 保证并发安全
}

// NewThriftManager 创建管理器实例
func NewThriftManager() *ThriftManager {
	return &ThriftManager{
		transports: make([]*Transport, 0),
	}
}

// Register 注册一个Thrift服务（只记录配置，不立即初始化连接）
// 参数：服务名、地址、客户端构造函数（用于后续初始化）
func (m *ThriftManager) Register(serviceName, addr string, factory func(thrift.TTransport) interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 先创建transport（不打开连接，等待Init时统一打开）
	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		panic("创建" + serviceName + " transport失败: " + err.Error())
	}

	// 构造客户端实例（基于transport）
	client := factory(transport)

	m.transports = append(m.transports, &Transport{
		ServiceName: serviceName,
		Addr:        addr,
		Transport:   transport,
		Client:      client,
	})
}

// Init 批量初始化所有注册的服务连接（打开transport）
func (m *ThriftManager) Init() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.transports {
		if err := t.Transport.Open(); err != nil {
			return err
		}
		log.Println("Thrift 服务已连接:", t.ServiceName, "->", t.Addr)
	}
	return nil
}

// Close 批量关闭所有服务的连接
func (m *ThriftManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, t := range m.transports {
		if t.Transport.IsOpen() {
			if err := t.Transport.Close(); err != nil {
				lastErr = err // 记录最后一个错误
			}
		}
	}
	return lastErr
}*/
