package job

import (
	"time"
	"your_project/library/logger"
)

// ExampleMinuteTask 示例：每分钟执行的任务
func ExampleMinuteTask() {
	logger.Info("job", "Running example minute task...")
	// 在这里编写你的业务逻辑
	// 例如：清理过期数据、发送通知等
}

// ExampleHourlyTask 示例：每小时执行的任务
func ExampleHourlyTask() {
	logger.Info("job", "Running example hourly task...")
	// 在这里编写你的业务逻辑
	// 例如：统计数据、生成报表等
}

// ExampleDailyTask 示例：每天执行的任务
func ExampleDailyTask() {
	logger.Info("job", "Executing minute task: %s", time.Now().Format("15:04:05"))
	// 在这里编写你的业务逻辑
	// 例如：备份数据、发送日报等
}
