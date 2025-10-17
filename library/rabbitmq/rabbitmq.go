package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"your_project/library/logger"
)

var conn *amqp.Connection
var channel *amqp.Channel

// Config RabbitMQ 配置
type Config struct {
	URL      string // 连接地址 amqp://user:pass@localhost:5672/
	Exchange string // 交换机名称
	Queue    string // 队列名称
}

// Init 初始化 RabbitMQ 连接
func Init(config Config) error {
	var err error

	// 建立连接
	conn, err = amqp.Dial(config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// 创建通道
	channel, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	logger.Info("rabbitmq", "RabbitMQ initialized successfully")
	return nil
}

// DeclareQueue 声明队列
func DeclareQueue(name string, durable, autoDelete bool) error {
	_, err := channel.QueueDeclare(
		name,       // 队列名称
		durable,    // 持久化
		autoDelete, // 自动删除
		false,      // 排他性
		false,      // 不等待
		nil,        // 参数
	)
	if err != nil {
		logger.Error("rabbitmq", "Failed to declare queue: %v", err)
		return err
	}
	logger.Info("rabbitmq", "Queue declared: %s", name)
	return nil
}

// DeclareExchange 声明交换机
func DeclareExchange(name, kind string, durable bool) error {
	err := channel.ExchangeDeclare(
		name,    // 交换机名称
		kind,    // 类型：direct, fanout, topic, headers
		durable, // 持久化
		false,   // 自动删除
		false,   // 内部
		false,   // 不等待
		nil,     // 参数
	)
	if err != nil {
		logger.Error("rabbitmq", "Failed to declare exchange: %v", err)
		return err
	}
	logger.Info("rabbitmq", "Exchange declared: %s (%s)", name, kind)
	return nil
}

// BindQueue 绑定队列到交换机
func BindQueue(queueName, routingKey, exchangeName string) error {
	err := channel.QueueBind(
		queueName,    // 队列名称
		routingKey,   // 路由键
		exchangeName, // 交换机名称
		false,        // 不等待
		nil,          // 参数
	)
	if err != nil {
		logger.Error("rabbitmq", "Failed to bind queue: %v", err)
		return err
	}
	logger.Info("rabbitmq", "Queue bound: %s -> %s (key: %s)", queueName, exchangeName, routingKey)
	return nil
}

// Publish 发布消息
func Publish(exchange, routingKey string, body []byte) error {
	err := channel.Publish(
		exchange,   // 交换机
		routingKey, // 路由键
		false,      // 强制
		false,      // 立即
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		logger.Error("rabbitmq", "Failed to publish message: %v", err)
		return err
	}
	logger.Debug("rabbitmq", "Message published to %s/%s", exchange, routingKey)
	return nil
}

// PublishToQueue 直接发布到队列
func PublishToQueue(queueName string, body []byte) error {
	return Publish("", queueName, body)
}

// Consume 消费消息
func Consume(queueName string, autoAck bool, handler func([]byte) error) error {
	msgs, err := channel.Consume(
		queueName, // 队列名称
		"",        // 消费者标签
		autoAck,   // 自动确认
		false,     // 排他性
		false,     // 不本地
		false,     // 不等待
		nil,       // 参数
	)
	if err != nil {
		logger.Error("rabbitmq", "Failed to register consumer: %v", err)
		return err
	}

	logger.Info("rabbitmq", "Started consuming from queue: %s", queueName)

	for msg := range msgs {
		logger.Debug("rabbitmq", "Received message: %s", string(msg.Body))

		// 处理消息
		if err := handler(msg.Body); err != nil {
			logger.Error("rabbitmq", "Failed to handle message: %v", err)
			if !autoAck {
				msg.Nack(false, true) // 拒绝并重新入队
			}
			continue
		}

		// 手动确认
		if !autoAck {
			msg.Ack(false)
		}
	}

	return nil
}

// Close 关闭连接
func Close() {
	if channel != nil {
		channel.Close()
	}
	if conn != nil {
		conn.Close()
	}
	logger.Info("rabbitmq", "RabbitMQ connection closed")
}

// ========================================
// 使用示例
// ========================================
/*
// 初始化
rabbitmq.Init(rabbitmq.Config{
    URL:      "amqp://guest:guest@localhost:5672/",
    Exchange: "my-exchange",
    Queue:    "my-queue",
})

// 声明队列
rabbitmq.DeclareQueue("my-queue", true, false)

// 声明交换机
rabbitmq.DeclareExchange("my-exchange", "direct", true)

// 绑定队列
rabbitmq.BindQueue("my-queue", "my-routing-key", "my-exchange")

// 发布消息
rabbitmq.Publish("my-exchange", "my-routing-key", []byte("Hello RabbitMQ"))

// 或直接发布到队列
rabbitmq.PublishToQueue("my-queue", []byte("Hello Queue"))

// 消费消息
rabbitmq.Consume("my-queue", false, func(body []byte) error {
    logger.Info("rabbitmq", "Processing: %s", string(body))
    return nil
})

// 关闭连接
defer rabbitmq.Close()
*/
