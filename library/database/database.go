package database

import (
	"log"
	"os"
	"sync"
	"time"
	"your_project/library/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var gormOnce sync.Once
var gormEngine *gorm.DB

// models 存储所有需要自动迁移的模型
var models []interface{}

// RegisterModels 注册需要自动迁移的模型
func RegisterModels(models ...interface{}) {
	if len(models) > 0 {
		for _, model := range models {
			registerModel(model)
		}
	}
}

// registerModel 注册单个模型
func registerModel(model interface{}) {
	models = append(models, model)
}

// Initial 初始化数据库连接
func Initial() {
	NewEngine()
	// 如果配置了自动迁移，则执行迁移
	if config.GetConfig().MysqlMaster.AutoMigrate {
		AutoMigrate()
	}
}

// NewEngine 获取数据库实例（单例模式）
func NewEngine() *gorm.DB {
	gormOnce.Do(func() {
		db, err := gorm.Open(mysql.Open(config.GetConfig().MysqlMaster.Dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 关闭外键约束
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.GetConfig().MysqlMaster.TablePrefix, // 表名前缀
				SingularTable: true,                                       // 使用单数表名
			},
		})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}

		// 配置日志级别
		logLevel := logger.Error
		if config.GetConfig().Server.Debug {
			logLevel = logger.Info
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢查询阈值
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)
		db.Logger = newLogger

		// 配置连接池
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get sql.DB: " + err.Error())
		}
		sqlDB.SetMaxIdleConns(config.GetConfig().MysqlMaster.MaxIdleConn)
		sqlDB.SetMaxOpenConns(config.GetConfig().MysqlMaster.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour)

		gormEngine = db
		log.Println("Database connected successfully")
	})
	return gormEngine
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() {
	if gormEngine == nil {
		log.Println("Database engine is not initialized, skipping AutoMigrate")
		return
	}

	if len(models) == 0 {
		log.Println("No models registered for AutoMigrate")
		return
	}

	log.Println("Starting database AutoMigrate...")
	err := gormEngine.AutoMigrate(models...)
	if err != nil {
		log.Printf("AutoMigrate failed: %v", err)
		panic("failed to migrate database: " + err.Error())
	}
	log.Printf("AutoMigrate completed successfully, migrated %d models", len(models))
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return gormEngine
}
