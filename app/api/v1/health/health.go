package health

import (
	"time"
	"your_project/app/view"
	"your_project/library/logger"

	"github.com/gofiber/fiber/v2"
)

// InitialHealthRoutes 注册健康检查路由
func InitialHealthRoutes(app *fiber.App) {
	app.Get("/v1/health", Health)
}

// ========================================
// 接口处理函数
// ========================================

// Health 健康检查
// @Summary      健康检查
// @Description  检查服务是否正常运行
// @Tags         健康检查
// @Accept       json
// @Produce      json
// @Success      200  {object}  view.Result{data=object{status=string,timestamp=int64,message=string}}  "成功"
// @Router       /v1/health [get]
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
