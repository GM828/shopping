package clients

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	nacosclient "github.com/kitex-contrib/config-nacos/client"
	"github.com/kitex-contrib/config-nacos/nacos"
	"github.com/kitex-contrib/config-nacos/utils"
	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	product "shopping/product_service/server/kitex_gen/product/productservice"
	"time"
)

var ProductClient product.Client

// configLog 用于打印 Nacos 配置内容
type configLog struct{}

// Apply 方法会在 Nacos 配置发生变化时被调用
func (cl *configLog) Apply(opt *utils.Options) {
	fn := func(cp *vo.ConfigParam) {
		klog.Infof("nacos config %v", cp)
	}
	// 将回调函数注册到 Nacos 自定义函数列表
	opt.NacosCustomFunctions = append(opt.NacosCustomFunctions, fn)
}

func Init() {
	// 创建 Nacos 解析器，用于服务发现
	r, err := resolver.NewDefaultNacosResolver()
	if err != nil {
		panic(err)
	}
	// 创建 Nacos 客户端，用于配置管理
	nacosClient, err := nacos.NewClient(nacos.Options{})
	if err != nil {
		panic(err)
	}

	cl := &configLog{}
	// 服务端服务名
	serviceName := "Product"
	// 客户端服务名
	clientName := "ProductClient"

	// 连接池配置
	poolConfig := connpool.IdleConfig{
		MaxIdlePerAddress: 10,
		MaxIdleGlobal:     100,
		MaxIdleTimeout:    time.Minute,
		MinIdlePerAddress: 2,
	}
	// 优雅配置
	opts := []client.Option{
		client.WithTransportProtocol(transport.TTHeader),                                 // 传输协议
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),                          // 元数据处理
		client.WithLongConnection(poolConfig),                                            // 连接池配置
		client.WithResolver(r),                                                           // 使用 Nacos 解析器进行服务发现
		client.WithSuite(nacosclient.NewSuite(serviceName, clientName, nacosClient, cl)), // 从服务中心拉取当前服务的相关配置
		client.WithRPCTimeout(2 * time.Second),                                           // 设置 RPC 超时时间
	}
	// 创建客户端
	c, err := product.NewClient(serviceName, opts...)
	if err != nil {
		log.Fatal(err)
	}
	ProductClient = c
	log.Println("ProductClient 初始化成功")
}
