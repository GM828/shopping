package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"log"
	"shopping/kitex_client/clients"
	"shopping/product_service/server/kitex_gen/product"
)

func main() {
	clients.Init()
	productClient := clients.ProductClient
	// 在这里使用 productClient 调用远程方法
	id := int64(1)
	resp, err := productClient.GetProductById(context.Background(), &product.ProductRequest{Id: &id})
	// 处理响应和错误
	if err != nil {
		// 尝试转换为业务异常
		bizErr, isBizErr := kerrors.FromBizStatusError(err)
		if isBizErr {
			// 处理业务异常
			log.Printf("[业务异常] Code=%d, Msg=%s", bizErr.BizStatusCode(), bizErr.BizMessage())
		} else {
			// 处理连接/系统异常
			log.Printf("[系统异常] %v", err)
			// 可根据错误类型细分处理:
			if kerrors.IsTimeoutError(err) {
				log.Println("-> 超时错误")
			} else {
				log.Println("-> 网络或其他系统错误")
			}
		}
		return
	}
	log.Println(resp)
}
