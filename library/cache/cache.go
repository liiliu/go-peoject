package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"your_project/library/common"
	"your_project/library/config"
	"your_project/library/logger"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client
var redisPrefix string

// Initial 初始化Redis连接
func Initial() {
	logger.Info("cache", "Initializing Redis...")

	redisPrefix = config.GetConfig().Redis.KeyPrefix

	redisClient = redis.NewClient(&redis.Options{
		Addr:        config.GetConfig().Redis.Host,
		Password:    config.GetConfig().Redis.Password,
		DB:          config.GetConfig().Redis.Db,
		IdleTimeout: config.GetConfig().Redis.IdleTimeout * time.Second,
		PoolSize:    config.GetConfig().Redis.PoolSize,
		MaxConnAge:  config.GetConfig().Redis.MaxConnAge * time.Second,
	})

	// 测试连接
	ctx1, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := redisClient.Ping(ctx1).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "Redis connection failed: %v", err)
		return
	}

	logger.Info("cache", "Redis connection initialized successfully: %s", pong)
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return redisClient
}

// Save 保存数据（JSON序列化）
func Save(name string, values interface{}, timeout time.Duration) error {
	serialized, err := json.Marshal(values)
	if err != nil {
		logger.Error(common.LogTagRedisError, "Marshal error: %v", err)
		return err
	}
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err = redisClient.Set(ctx, key, serialized, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "Set error: %v", err)
		return err
	}
	return nil
}

// Get 获取数据（JSON反序列化）
func Get(name string, values interface{}) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("key not found")
		}
		logger.Error(common.LogTagRedisError, "Get error: %v", err)
		return err
	}

	err = json.Unmarshal(value, &values)
	if err != nil {
		logger.Error(common.LogTagRedisError, "Unmarshal error: %v", err)
		return err
	}

	return nil
}

// SaveString 保存字符串
func SaveString(name string, str string, timeout time.Duration) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Set(ctx, key, str, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "SaveString error: %v", err)
		return err
	}
	return nil
}

// GetString 获取字符串
func GetString(name string) string {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.Error(common.LogTagRedisError, "GetString error: %v", err)
		}
		return ""
	}
	return value
}

// Remove 删除缓存
func Remove(name string) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Del(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "Remove error: %v", err)
		return err
	}
	return nil
}

// Expire 设置过期时间
func Expire(name string, timeout time.Duration) error {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	_, err := redisClient.Expire(ctx, key, timeout).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "Expire error: %v", err)
		return err
	}
	return nil
}

// Exists 检查键是否存在
func Exists(name string) bool {
	key := fmt.Sprintf("%s:%s", redisPrefix, name)
	val, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		logger.Error(common.LogTagRedisError, "Exists error: %v", err)
		return false
	}
	return val > 0
}
