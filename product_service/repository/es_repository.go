package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
	"shopping/conn"
	"shopping/product_service/model"
)

type EsRepository interface {
	EsIndexRepository
	EsDocRepository
}

// 定义索引相关的操作
type EsIndexRepository interface {
	// 创建索引
	CreateIndex(mapping string) error
	// 判断索引是否存在
	ExistsIndex() bool
	// 删除索引
	DeleteIndex() bool
	// 添加索引字段
	AddIndexField(mapping string) bool
}

// 定义文档相关的操作
type EsDocRepository interface {
	// 创建文档
	CreateDoc(mo *model.ProductFullMO) error
	// 批量创建文档
	BulkCreateDocs(mos []*model.ProductFullMO) error
	// 删除文档（通过业务id）
	DeleteDoc(id int64) error
	// 模糊查询文档
	SearchDocs(query string) ([]*model.ProductFullMO, error)
}

type EsRepositoryImpl struct {
	esClient  *elastic.Client
	indexName string // 固定索引名，如 "product_index"
}

func NewEsRepository() EsRepository {
	return &EsRepositoryImpl{
		esClient:  conn.EsClient,
		indexName: "product_index",
	}
}

func (r *EsRepositoryImpl) CreateIndex(mapping string) error {
	// 先判断索引是否存在
	if r.ExistsIndex() {
		return errors.New("索引已存在")
	}
	// 创建索引
	createdIndex, err := r.esClient.CreateIndex(r.indexName).BodyString(mapping).Do(context.Background())
	if err != nil {
		return errors.New("CreateIndex error: " + err.Error())
	}
	log.Println("索引创建成功：" + createdIndex.Index)
	return nil
}

func (r *EsRepositoryImpl) ExistsIndex() bool {
	exists, err := r.esClient.IndexExists(r.indexName).Do(context.Background())
	if err != nil {
		log.Println("ExistsIndex error:", err)
		return false
	}
	if exists {
		log.Println("索引已存在")
	} else {
		log.Println("索引不存在")
	}
	return exists
}

func (r *EsRepositoryImpl) DeleteIndex() bool {
	deletedIndex, err := r.esClient.DeleteIndex(r.indexName).Do(context.Background())
	if err != nil {
		log.Println("DeleteIndex error:", err)
		return false
	}
	return deletedIndex.Acknowledged
}

func (r *EsRepositoryImpl) AddIndexField(mapping string) bool {
	putResp, err := r.esClient.PutMapping().
		Index(r.indexName).
		BodyString(mapping).
		Do(context.Background())

	if err != nil {
		log.Println("AddIndexField error:", err)
		return false
	}
	return putResp.Acknowledged
}

func (r *EsRepositoryImpl) CreateDoc(mo *model.ProductFullMO) error {
	resp, err := r.esClient.Index().
		Index(r.indexName).
		BodyJson(mo).
		Do(context.Background())

	if err != nil {
		return errors.New("CreateDoc error: " + err.Error())
	}

	log.Println("CreateDoc success, Id: ", resp.Id)
	return nil
}

func (r *EsRepositoryImpl) BulkCreateDocs(mos []*model.ProductFullMO) error {
	bulkRequest := r.esClient.Bulk().Refresh("true")
	for _, mo := range mos {
		req := elastic.NewBulkIndexRequest().
			Index(r.indexName).
			Doc(mo)
		bulkRequest = bulkRequest.Add(req)
	}
	bulkResp, err := bulkRequest.Do(context.Background())
	if err != nil {
		return errors.New("BulkCreateDocs error: " + err.Error())
	}
	log.Println("BulkCreateDocs success, Items count: ", len(bulkResp.Items))
	return nil
}

func (r *EsRepositoryImpl) DeleteDoc(id int64) error {
	// 删除条件
	query := elastic.NewTermQuery("id", id)

	resp, err := conn.EsClient.DeleteByQuery().
		Index(r.indexName).
		Query(query).
		Refresh("true"). // 立即刷新索引，使删除操作生效
		Do(context.Background())

	if err != nil {
		return errors.New("DeleteDocByBusinessId error: " + err.Error())
	}

	log.Println("DeleteDocByBusinessId success, Deleted count:", resp.Deleted)
	return nil
}

func (r *EsRepositoryImpl) SearchDocs(query string) ([]*model.ProductFullMO, error) {
	// 可以支持模糊查询的字段有两个：productName 和 description
	// 先根据 productName 查询
	matchQuery := elastic.NewMatchQuery("productName", query)
	searchResult, err := r.esClient.Search().
		Index(r.indexName).
		Query(matchQuery).
		From(0).Size(100).
		Do(context.Background())
	if err != nil {
		log.Println("SearchDocs error:", err)
		return nil, err
	}
	// 如果没有结果，再根据 description 查询
	if searchResult.TotalHits() == 0 {
		matchQuery = elastic.NewMatchQuery("descriptions", query)
		searchResult, err = r.esClient.Search().
			Index(r.indexName).
			Query(matchQuery).
			From(0).Size(100).
			Do(context.Background())
		if err != nil {
			log.Println("SearchDocs error:", err)
			return nil, err
		}
	}
	// 类型转换
	var results []*model.ProductFullMO
	for _, hit := range searchResult.Hits.Hits {
		var d model.ProductFullMO
		err := json.Unmarshal(hit.Source, &d)
		if err != nil {
			return nil, errors.New("SearchDocs unmarshal error: " + err.Error())
		}
		results = append(results, &d)
	}
	return results, nil
}
