package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"your_project/library/logger"
)

var client *elasticsearch.Client

// Config Elasticsearch 配置
type Config struct {
	Addresses []string // ES 服务器地址列表
	Username  string   // 用户名
	Password  string   // 密码
}

// Init 初始化 Elasticsearch 客户端
func Init(config Config) error {
	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
		Username:  config.Username,
		Password:  config.Password,
	}

	var err error
	client, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create ES client: %v", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		return fmt.Errorf("failed to get ES info: %v", err)
	}
	defer res.Body.Close()

	logger.Info("elasticsearch", "Elasticsearch client initialized successfully")
	return nil
}

// CreateIndex 创建索引
func CreateIndex(indexName string, mapping map[string]interface{}) error {
	body, _ := json.Marshal(mapping)

	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	logger.Info("elasticsearch", "Index created: %s", indexName)
	return nil
}

// IndexDocument 索引文档
func IndexDocument(indexName, docID string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	res, err := client.Index(
		indexName,
		bytes.NewReader(body),
		client.Index.WithDocumentID(docID),
		client.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document: %s", res.String())
	}

	logger.Debug("elasticsearch", "Document indexed: %s/%s", indexName, docID)
	return nil
}

// GetDocument 获取文档
func GetDocument(indexName, docID string) (map[string]interface{}, error) {
	res, err := client.Get(indexName, docID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("document not found: %s", docID)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// Search 搜索文档
func Search(indexName string, query map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(indexName),
		client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search failed: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// 解析搜索结果
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	docs := make([]map[string]interface{}, len(hits))

	for i, hit := range hits {
		h := hit.(map[string]interface{})
		docs[i] = h["_source"].(map[string]interface{})
	}

	return docs, nil
}

// DeleteDocument 删除文档
func DeleteDocument(indexName, docID string) error {
	res, err := client.Delete(indexName, docID)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to delete document: %s", res.String())
	}

	logger.Info("elasticsearch", "Document deleted: %s/%s", indexName, docID)
	return nil
}

// DeleteIndex 删除索引
func DeleteIndex(indexName string) error {
	res, err := client.Indices.Delete([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to delete index: %s", res.String())
	}

	logger.Info("elasticsearch", "Index deleted: %s", indexName)
	return nil
}

// BulkIndex 批量索引
func BulkIndex(indexName string, docs []interface{}) error {
	var buf bytes.Buffer

	for _, doc := range docs {
		// 元数据
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
			},
		}
		metaJSON, _ := json.Marshal(meta)
		buf.Write(metaJSON)
		buf.WriteByte('\n')

		// 文档数据
		docJSON, _ := json.Marshal(doc)
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	res, err := client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index failed: %s", res.String())
	}

	logger.Info("elasticsearch", "Bulk indexed %d documents to %s", len(docs), indexName)
	return nil
}

// ========================================
// 使用示例
// ========================================
/*
// 初始化
elasticsearch.Init(elasticsearch.Config{
    Addresses: []string{"http://localhost:9200"},
    Username:  "elastic",
    Password:  "password",
})

// 创建索引
mapping := map[string]interface{}{
    "mappings": map[string]interface{}{
        "properties": map[string]interface{}{
            "title": map[string]interface{}{"type": "text"},
            "content": map[string]interface{}{"type": "text"},
            "created_at": map[string]interface{}{"type": "date"},
        },
    },
}
elasticsearch.CreateIndex("articles", mapping)

// 索引文档
doc := map[string]interface{}{
    "title":      "Elasticsearch Tutorial",
    "content":    "This is a tutorial about Elasticsearch",
    "created_at": time.Now(),
}
elasticsearch.IndexDocument("articles", "1", doc)

// 搜索
query := map[string]interface{}{
    "query": map[string]interface{}{
        "match": map[string]interface{}{
            "content": "Elasticsearch",
        },
    },
}
results, _ := elasticsearch.Search("articles", query)

// 获取文档
doc, _ := elasticsearch.GetDocument("articles", "1")

// 删除文档
elasticsearch.DeleteDocument("articles", "1")
*/
