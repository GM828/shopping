package main

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"shopping/database"
	"shopping/thrift_def/gen-go/user"
	"shopping/user_service"
)

func main() {
	// 初始化数据库连接
	database.Init()

	// 1. 获取 IDL 当中服务的实现类
	userHandler, err := user_service.InitializeUserHandler()
	if err != nil {
		panic(err)
	}
	// 2. 创建 Thrift 处理器(绑定实现类)
	processor := user.NewUserServiceProcessor(userHandler)

	// 3. 配置传输层和协议层
	// 创建服务端传输层,监听 9090 端口
	serverTransport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		panic(err)
	}
	// 创建传输层工厂(带缓冲,8KB)
	// transportFactory := thrift.NewTBufferedTransportFactory(8192)
	// 优化：长连接使用 TFramedTransportFactory
	transportFactory := thrift.NewTFramedTransportFactory(
		thrift.NewTBufferedTransportFactory(8192),
	)
	// 创建协议层工厂(二进制协议)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	// 4. 创建 goroutine pool
	maxWorkers := 100 // 最大并发处理数(worker 数量)
	// 创建任务队列,容量为 worker 数量的 2 倍
	jobQueue := make(chan thrift.TTransport, maxWorkers*2)

	// 启动固定数量的 worker goroutine
	for i := 0; i < maxWorkers; i++ {
		go worker(processor, transportFactory, protocolFactory, jobQueue)
	}

	println("Thrift 服务端启动,监听 :9090...")
	// 开始监听端口
	serverTransport.Listen()

	// 主循环:接受客户端连接并分发到 worker pool
	for {
		// 接受一个新的客户端连接
		client, err := serverTransport.Accept()
		// 日志：显示接受到的连接标识
		log.Printf("[服务端] 接受连接 | Transport地址: %p | 远程地址: %s | 本地地址: %s\n",
			client,
			client.(*thrift.TSocket).Conn().RemoteAddr().String(),
			client.(*thrift.TSocket).Conn().LocalAddr().String())

		if err != nil {
			continue
		}
		// 将连接发送到任务队列,由空闲的 worker 处理
		jobQueue <- client
	}
}

// worker 函数:从任务队列中获取连接并处理请求
func worker(processor thrift.TProcessor, transportFactory thrift.TTransportFactory, protocolFactory thrift.TProtocolFactory, jobs <-chan thrift.TTransport) {
	// 持续从任务队列中获取客户端连接
	for client := range jobs {
		// 为当前连接创建带缓冲的传输层
		trans, err := transportFactory.GetTransport(client)
		if err != nil {
			client.Close()
			continue
		}
		// 创建输入协议(用于反序列化请求)
		inputProtocol := protocolFactory.GetProtocol(trans)
		// 创建输出协议(用于序列化响应)
		outputProtocol := protocolFactory.GetProtocol(trans)

		// 循环处理该连接上的多个请求(支持连接复用)
		for {
			// 处理一个 RPC 请求
			// keepOpen: 客户端是否希望保持连接
			keepOpen, err := processor.Process(context.Background(), inputProtocol, outputProtocol)
			// 如果发生错误或客户端要求关闭连接,则退出循环
			if !keepOpen {
				log.Println("[服务端] 客户端关闭连接")
				break
			}
			// 关键修改:只在严重错误时退出循环
			if err != nil && shouldCloseConnection(err) {
				log.Printf("[服务端] 严重错误,关闭连接: %v\n", err)
				break
			}
		}
		// 关闭传输层连接
		trans.Close()
	}
}

// 黑名单模式：明确哪些错误必须断开连接
func shouldCloseConnection(err error) bool {
	// 1. 连接层错误（网络断开等）→ 必须关闭
	if _, ok := err.(thrift.TTransportException); ok {
		return true
	}

	// 2. 协议层错误（数据格式错误）→ 必须关闭
	if _, ok := err.(thrift.TProtocolException); ok {
		return true
	}

	// 3. 其他错误（TApplicationException、业务异常）→ 不关闭
	return false
}

/*// 普通的单线程服务器
func main() {
	// 初始化数据库
	database.Init()

	// 1. 获取IDL当中服务的实现类
	userHandler, err := user_service.InitializeUserHandler()
	if err != nil {
		panic(err)
	}
	// 2. 创建 Thrift 处理器（绑定实现类）
	processor := user.NewUserServiceProcessor(userHandler)

	// 3. 配置传输层和协议层
	transport, err := thrift.NewTServerSocket(":9090") // 监听 9090 端口
	if err != nil {
		panic(err)
	}
	transportFactory := thrift.NewTBufferedTransportFactory(8192) // 带缓冲，缓冲区大小为8kb
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()  // 二进制协议

	// 4. 创建并启动服务器
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	println("Thrift 服务端启动，监听 :9090...")
	if err := server.Serve(); err != nil {
		panic(err)
	}
}*/

/*// 每个连接一个 goroutine
func main() {
	// 初始化数据库连接
	database.Init()

	// 1. 获取 IDL 当中服务的实现类
	userHandler, err := user_service.InitializeUserHandler()
	if err != nil {
		panic(err)
	}
	// 2. 创建 Thrift 处理器(绑定实现类)
	processor := user.NewUserServiceProcessor(userHandler)

	// 3. 创建服务端传输层,监听 9090 端口
	serverTransport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		panic(err)
	}

	println("Thrift 服务端启动,监听 :9090...")
	// 开始监听端口
	serverTransport.Listen()

	// 主循环:接受客户端连接
	for {
		// 接受一个新的客户端连接
		client, err := serverTransport.Accept()
		if err != nil {
			continue
		}
		// 为每个连接启动一个 goroutine 进行处理(无限制并发)
		go func(c thrift.TTransport) {
			// 4. 为当前连接创建传输层工厂(带缓冲,8KB)
			transportFactory := thrift.NewTBufferedTransportFactory(8192)
			// 5. 创建协议层工厂(二进制协议)
			protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

			// 获取带缓冲的传输层
			trans, err := transportFactory.GetTransport(c)
			if err != nil {
				c.Close()
				return
			}

			// 创建输入和输出协议(用于序列化/反序列化)
			inputProtocol := protocolFactory.GetProtocol(trans)
			outputProtocol := protocolFactory.GetProtocol(trans)

			// 循环处理该连接上的多个请求
			for {
				// 处理一个 RPC 请求
				// keepOpen: 是否保持连接
				keepOpen, err := processor.Process(context.Background(), inputProtocol, outputProtocol)
				// 如果发生错误或客户端要求关闭连接,则退出循环
				if err != nil || !keepOpen {
					break
				}
			}
			// 关闭传输层连接
			trans.Close()
		}(client)
	}
}*/
