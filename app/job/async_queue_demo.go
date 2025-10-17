package job

import (
	"encoding/json"
	"your_project/library/logger"
	"your_project/library/queue"
)

// QueueTaskDemo 队列异步任务示例
// 这是一个基于 worker pool 的异步任务处理模板
// 可以根据实际业务需求改造此模板
func QueueTaskDemo() {
	logger.Info("job", "Queue task demo started")
	
	defer func() {
		if err := recover(); err != nil {
			logger.Error("job", "Queue task panic: %v", err)
		}
	}()

	// 初始化 worker pool（协程池）
	// 参数：并发 worker 数量
	workerPool := queue.NewWorkerPool(3)
	workerPool.Run()

	// 模拟从队列中读取任务
	// 实际使用时，可以从 Redis、RabbitMQ、Kafka 等消息队列读取
	// 示例：使用 Redis BLPOP 阻塞读取队列
	/*
		queueName := "YOUR_QUEUE_NAME"
		for {
			// 从 Redis 队列中读取任务数据
			reply, err := cache.BLPOP(queueName, 0)
			if err != nil {
				logger.Error("job", "Failed to read from queue: %v", err)
				time.Sleep(30 * time.Second)
				continue
			}

			// 解析队列数据
			taskData := new(YourTaskStruct)
			err = json.Unmarshal([]byte(reply), &taskData)
			if err != nil {
				logger.Error("job", "Failed to unmarshal task data: %v", err)
				continue
			}

			// 创建任务实例并提交到 worker pool
			task := &DemoTask{Data: taskData}
			workerPool.JobQueue <- task
		}
	*/

	// 当前为示例代码，实际使用时请替换为真实的队列读取逻辑
	logger.Info("job", "Queue task demo: waiting for tasks from queue...")
	
	// 保持任务运行
	select {}
}

// DemoTaskData 任务数据结构示例
// 根据实际业务定义您的任务数据结构
type DemoTaskData struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
	// 添加更多字段...
}

// DemoTask 实现 queue.Job 接口的任务
type DemoTask struct {
	Data *DemoTaskData
}

// Do 实现 Job 接口，执行具体的任务逻辑
func (t *DemoTask) Do() {
	// 任务数据验证
	if t.Data == nil {
		logger.Error("job", "Task data is nil")
		return
	}

	logger.Info("job", "Processing task: ID=%d, Content=%s", t.Data.ID, t.Data.Content)

	// ========================================
	// 在这里编写您的业务逻辑
	// ========================================
	// 示例场景：
	// 1. 发送短信/邮件
	// 2. 处理文件上传
	// 3. 图片/视频处理
	// 4. 数据导入/导出
	// 5. 调用第三方 API
	// 6. 数据统计/计算
	// ========================================

	/*
		// 示例：查询数据库
		var record YourModel
		if err := database.NewEngine().Where("id = ?", t.Data.ID).First(&record).Error; err != nil {
			logger.Error("job", "Failed to query database: %v", err)
			return
		}

		// 示例：执行业务逻辑
		result, err := YourBusinessLogic(record)
		if err != nil {
			logger.Error("job", "Business logic failed: %v", err)
			// 更新状态为失败
			record.Status = StatusFailed
		} else {
			logger.Info("job", "Business logic succeeded: %v", result)
			// 更新状态为成功
			record.Status = StatusSuccess
		}

		// 示例：保存结果
		if err := database.NewEngine().Save(&record).Error; err != nil {
			logger.Error("job", "Failed to save result: %v", err)
			return
		}
	*/

	logger.Info("job", "Task completed: ID=%d", t.Data.ID)
}

// ========================================
// 辅助函数示例
// ========================================

// parseTaskData 解析任务数据的辅助函数
func parseTaskData(data []byte) (*DemoTaskData, error) {
	taskData := new(DemoTaskData)
	err := json.Unmarshal(data, taskData)
	if err != nil {
		return nil, err
	}
	return taskData, nil
}

// validateTaskData 验证任务数据
func validateTaskData(data *DemoTaskData) bool {
	if data == nil {
		return false
	}
	if data.ID == 0 {
		return false
	}
	// 添加更多验证逻辑...
	return true
}

// ========================================
// 使用说明
// ========================================
/*
使用此模板创建异步队列任务：

1. 定义任务数据结构
   - 修改 DemoTaskData 结构体，添加您的字段

2. 实现任务逻辑
   - 在 DemoTask.Do() 方法中编写业务逻辑

3. 接入消息队列
   - 在 QueueTaskDemo() 中实现队列读取逻辑
   - 支持 Redis、RabbitMQ、Kafka 等

4. 启动任务
   - 在 commands/worker/worker.go 中调用 QueueTaskDemo()
   - 或在 job.go 中添加启动逻辑

5. 调整 worker pool 大小
   - 根据任务量和服务器资源调整并发数
   - workerPool := queue.NewWorkerPool(10)

示例：推送任务到 Redis 队列
```go
// 在其他地方推送任务到队列
taskData := &DemoTaskData{
    ID:      1,
    Content: "example task",
}
data, _ := json.Marshal(taskData)
cache.RPUSH("YOUR_QUEUE_NAME", string(data))
```
*/
