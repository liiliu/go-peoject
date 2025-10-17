package config

import (
	"log"
	"time"

	"github.com/koding/multiconfig"
)

var c *Config

// Init 初始化配置
func Init() {
	log.Println("Initializing config...")
	m := multiconfig.NewWithPath("config/config.toml")
	serverConf := new(Config)
	m.MustLoad(serverConf)

	c = serverConf
	log.Printf("Config loaded: %+v", c.Server)
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return c
}

// Config 配置结构体
type Config struct {
	Server      ServerConfig
	MysqlMaster MysqlMasterConfig
	Redis       RedisConfig
	JWT         JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Application   string
	Name          string
	Port          int64
	Debug         bool
	Env           string
	LogFileFolder string
	TmpPath       string
}

// MysqlMasterConfig MySQL主库配置
type MysqlMasterConfig struct {
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
	TablePrefix string
	AutoMigrate bool
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host        string
	Password    string
	Db          int
	IdleTimeout time.Duration
	PoolSize    int
	MaxConnAge  time.Duration
	KeyPrefix   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret         string
	ExpireTime     time.Duration
	MaxRefreshTime time.Duration
}
