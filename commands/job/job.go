package job

import (
	"your_project/app/job"
	"your_project/app/model"
	"your_project/library/cache"
	"your_project/library/config"
	"your_project/library/database"
	"your_project/library/logger"
)

// Initial 初始化所有组件
func Initial() {
	config.Init()
	logger.Init()
	// 注册需要自动迁移的模型
	registerModels()
	database.Initial()
	cache.Initial()
}

// registerModels 注册所有需要自动迁移的数据库模型
func registerModels() {
	database.RegisterModels(
		&model.User{},
		// 在这里添加其他模型...
	)
}

// Run 启动定时任务服务
func Run() {
	Initial()
	logger.Info("job", "Starting job scheduler...")
	job.Run()
}
