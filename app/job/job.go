package job

import (
	"your_project/library/logger"

	"github.com/robfig/cron"
)

var cronEngine *cron.Cron

// Run 运行 Job 服务（定时任务 + 异步队列任务）
func Run() {
	logger.Info("job", "Starting job service...")
	
	// 启动异步队列任务（在单独的协程中运行）
	// go QueueTaskDemo()
	
	// 启动定时任务调度器（会阻塞）
	RunScheduler()
}

// RunScheduler 启动调度器
func RunScheduler() {
	logger.Info("job", "Starting job scheduler...")
	cronEngine = cron.New()
	scheduleJobs()
	cronEngine.Run()
}

// scheduleJobs 注册定时任务
func scheduleJobs() {
	// 示例1: 每分钟执行一次
	cronEngine.AddFunc("0 * * * * *", func() {
		ExampleMinuteTask()
	})

	// 示例2: 每小时执行一次
	cronEngine.AddFunc("0 0 * * * *", func() {
		ExampleHourlyTask()
	})

	// 示例3: 每天凌晨3点执行
	cronEngine.AddFunc("0 0 3 * * *", func() {
		ExampleDailyTask()
	})

	logger.Info("job", "Scheduled jobs registered")
}
