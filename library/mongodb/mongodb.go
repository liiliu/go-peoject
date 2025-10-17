package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"your_project/library/logger"
)

var client *mongo.Client
var database *mongo.Database

// Config MongoDB 配置
type Config struct {
	URI      string // 连接字符串
	Database string // 数据库名称
	Username string // 用户名（可选）
	Password string // 密码（可选）
}

// Init 初始化 MongoDB
func Init(config Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	if config.Username != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})
	}

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping: %v", err)
	}

	database = client.Database(config.Database)
	logger.Info("mongodb", "MongoDB connected: database=%s", config.Database)
	return nil
}

// GetCollection 获取集合
func GetCollection(name string, tag string) *mongo.Collection {
	return database.Collection(name)
}

// InsertOne 插入文档
func InsertOne(collectionName string, document interface{}, tag string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := GetCollection(collectionName, tag).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	logger.Debug("mongodb", "Document inserted: %v", result.InsertedID)
	return result.InsertedID, nil
}

// FindOne 查询单个文档
func FindOne(collectionName string, filter interface{}, result interface{}, tag string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return GetCollection(collectionName).FindOne(ctx, filter).Decode(result)
}

// UpdateOne 更新文档
func UpdateOne(collectionName string, filter, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetCollection(collectionName).UpdateOne(ctx, filter, update)
	return err
}

// DeleteOne 删除文档
func DeleteOne(collectionName string, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetCollection(collectionName).DeleteOne(ctx, filter)
	return err
}

// Close 关闭连接
func Close() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return client.Disconnect(ctx)
	}
	return nil
}
