package health

import (
	"time"
	"your_project/app/view"
	"your_project/library/logger"

	"github.com/gofiber/fiber/v2"
)

// Health 健康检查
func Health(c *fiber.Ctx) error {
	// 使用带 traceId 的日志
	logger.InfoWithTrace(c, "health", "健康检查请求")
	
	// 返回响应，自动带上 traceId
	return c.JSON(view.SuccessWithCtx(c, fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"message":   "Service is running",
	}))
}

// InitialHealthRoutes 注册健康检查路由
func InitialHealthRoutes(app *fiber.App) {
	app.Get("/v1/health", Health)
}
