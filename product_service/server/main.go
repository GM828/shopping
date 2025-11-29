package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/config-nacos/nacos"
	nacosserver "github.com/kitex-contrib/config-nacos/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"log"
	"net"
	"shopping/conn"
	"shopping/database"
	"shopping/product_service"
	product "shopping/product_service/server/kitex_gen/product/productservice"
)

func main() {
	// 初始化MySQL数据库连接
	database.Init()

	// 初始化ES数据库连接
	err := conn.EsConnect()
	if err != nil {
		log.Println("EsConnect error:", err)
	} else {
		log.Println("EsConnect success : " + conn.EsClient.String())
	}

	// 设置日志级别为 Debug，方便观察日志
	klog.SetLevel(klog.LevelDebug)

	// 创建 Nacos 注册中心实例，用于服务注册
	r, err := registry.NewDefaultNacosRegistry()
	if err != nil {
		panic(err)
	}

	// 创建 Nacos 客户端，用于拉取配置
	nacosClient, err := nacos.NewClient(nacos.Options{})
	if err != nil {
		panic(err)
	}

	// 该名称用于服务注册和配置中心的标识，客户端通过此名称查找服务
	serviceName := "Product"

	// 自定义端口，这里示例设置为9999，IP指定为0.0.0.0（允许所有网卡访问）
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal(err)
	}

	// 优雅配置
	opts := []server.Option{
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),                          // 使用 TTHeader 元数据处理器
		server.WithServiceAddr(addr),                                                     // 使用自定义地址
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}), // 设置服务基本信息（服务名必须与注册中心一致）
		server.WithRegistry(r),                                                           // 注册服务
		server.WithSuite(nacosserver.NewSuite(serviceName, nacosClient)),                 // 从服务中心拉取当前服务的相关配置
	}

	// 对handler进行依赖注入
	handler, err := product_service.InitializeProductHandler()
	if err != nil {
		log.Fatal("初始化 Handler 失败:", err)
	}

	// 创建服务器
	svr := product.NewServer(handler, opts...)
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
