package health

import (
	"time"
	"your_project/app/view"

	"github.com/gofiber/fiber/v2"
)

// Health 健康检查
func Health(c *fiber.Ctx) error {
	return c.JSON(view.SuccessResult(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"message":   "Service is running",
	}))
}

// InitialHealthRoutes 注册健康检查路由
func InitialHealthRoutes(app *fiber.App) {
	app.Get("/v1/health", Health)
}
