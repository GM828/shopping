package conn

import "github.com/olivere/elastic/v7"

// EsClient 全局 Elasticsearch 客户端
var EsClient *elastic.Client

func EsConnect() error {
	// 配置 Elasticsearch 客户端选项
	esopts := []elastic.ClientOptionFunc{
		elastic.SetURL("http://192.168.163.133:9200"), // Elasticsearch 服务器地址
		elastic.SetSniff(false),                       // 关闭节点嗅探功能
		elastic.SetBasicAuth("", ""),                  // 如果有用户名和密码，填写在这里
	}
	// 创建 Elasticsearch 客户端
	client, err := elastic.NewClient(esopts...)
	if err != nil {
		return err
	}
	// 连接成功，赋值
	EsClient = client
	return nil
}
