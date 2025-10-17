package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"your_project/library/logger"
)

var (
	writer *kafka.Writer
	reader *kafka.Reader
)

// Config Kafka 配置
type Config struct {
	Brokers []string // Kafka 服务器地址列表
	Topic   string   // 主题
	GroupID string   // 消费者组 ID
}

// InitProducer 初始化生产者
func InitProducer(brokers []string, topic string) {
	writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        false, // 同步发送
	}
	logger.Info("kafka", "Kafka producer initialized: topic=%s", topic)
}

// InitConsumer 初始化消费者
func InitConsumer(brokers []string, topic, groupID string) {
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	logger.Info("kafka", "Kafka consumer initialized: topic=%s, groupID=%s", topic, groupID)
}

// Produce 发送消息
func Produce(ctx context.Context, key, value []byte) error {
	if writer == nil {
		logger.Error("kafka", "Kafka producer not initialized")
		return nil
	}

	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})

	if err != nil {
		logger.Error("kafka", "Failed to produce message: %v", err)
		return err
	}

	logger.Debug("kafka", "Message produced: key=%s", string(key))
	return nil
}

// ProduceJSON 发送 JSON 消息
func ProduceJSON(ctx context.Context, key string, value interface{}) error {
	// TODO: 序列化为 JSON
	// jsonData, _ := json.Marshal(value)
	// return Produce(ctx, []byte(key), jsonData)
	return nil
}

// Consume 消费消息（阻塞）
func Consume(ctx context.Context, handler func(key, value []byte) error) error {
	if reader == nil {
		logger.Error("kafka", "Kafka consumer not initialized")
		return nil
	}

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("kafka", "Failed to read message: %v", err)
			continue
		}

		logger.Debug("kafka", "Message consumed: key=%s, offset=%d", string(msg.Key), msg.Offset)

		// 处理消息
		if err := handler(msg.Key, msg.Value); err != nil {
			logger.Error("kafka", "Failed to handle message: %v", err)
			// 根据需要决定是否继续消费
			continue
		}
	}
}

// Close 关闭连接
func Close() {
	if writer != nil {
		writer.Close()
		logger.Info("kafka", "Kafka producer closed")
	}
	if reader != nil {
		reader.Close()
	logger.Info("kafka", "Kafka consumer closed")
	}
}

// ========================================
// 使用示例
// ========================================
/*
// 初始化生产者
kafka.InitProducer([]string{"localhost:9092"}, "my-topic")

// 发送消息
ctx := context.Background()
kafka.Produce(ctx, []byte("key1"), []byte("message content"))

// 初始化消费者
kafka.InitConsumer([]string{"localhost:9092"}, "my-topic", "my-group")

// 消费消息
kafka.Consume(ctx, func(key, value []byte) error {
    logger.Info("kafka", "Received: %s = %s", string(key), string(value))
    return nil
})
*/
