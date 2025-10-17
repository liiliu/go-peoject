package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"your_project/library/logger"
)

var client mqtt.Client

// Config MQTT 配置
type Config struct {
	Broker   string // MQTT 服务器地址 (tcp://localhost:1883)
	ClientID string // 客户端 ID
	Username string // 用户名
	Password string // 密码
}

// Init 初始化 MQTT 客户端
func Init(config Config) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetAutoReconnect(true)

	// 连接回调
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		logger.Info("mqtt", "MQTT connected to %s", config.Broker)
	})

	// 连接丢失回调
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		logger.Error("mqtt", "MQTT connection lost: %v", err)
	})

	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	logger.Info("mqtt", "MQTT client initialized successfully")
	return nil
}

// Publish 发布消息
func Publish(topic string, payload interface{}, qos byte, retained bool) error {
	if client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := client.Publish(topic, qos, retained, payload)
	token.Wait()

	if token.Error() != nil {
		logger.Error("mqtt", "Failed to publish to %s: %v", topic, token.Error())
		return token.Error()
	}

	logger.Debug("mqtt", "Message published to %s", topic)
	return nil
}

// Subscribe 订阅主题
func Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT client not connected")
	}

	token := client.Subscribe(topic, qos, callback)
	token.Wait()

	if token.Error() != nil {
		logger.Error("mqtt", "Failed to subscribe to %s: %v", topic, token.Error())
		return token.Error()
	}

	logger.Info("mqtt", "Subscribed to topic: %s", topic)
	return nil
}

// Unsubscribe 取消订阅
func Unsubscribe(topics ...string) error {
	if client == nil {
		return fmt.Errorf("MQTT client not initialized")
	}

	token := client.Unsubscribe(topics...)
	token.Wait()

	if token.Error() != nil {
		return token.Error()
	}

	logger.Info("mqtt", "Unsubscribed from topics: %v", topics)
	return nil
}

// Disconnect 断开连接
func Disconnect() {
	if client != nil && client.IsConnected() {
		client.Disconnect(250)
		logger.Info("mqtt", "MQTT client disconnected")
	}
}

// IsConnected 检查连接状态
func IsConnected() bool {
	return client != nil && client.IsConnected()
}

// ========================================
// 使用示例
// ========================================
/*
// 初始化
mqtt.Init(mqtt.Config{
    Broker:   "tcp://localhost:1883",
    ClientID: "go-client-1",
    Username: "user",
    Password: "pass",
})

// 发布消息
mqtt.Publish("sensor/temperature", "25.5", 1, false)

// 订阅消息
mqtt.Subscribe("sensor/#", 1, func(client mqtt.Client, msg mqtt.Message) {
    logger.Info("mqtt", "Received: %s = %s", msg.Topic(), string(msg.Payload()))
})

// 断开连接
mqtt.Disconnect()
*/
