package router

import (
	"fmt"
	"your_project/app/api/middleware"
	"your_project/app/api/v1/auth"
	"your_project/app/api/v1/health"
	"your_project/library/config"
	"your_project/library/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// Run 启动路由
func Run() {
	logger.Info("api", "Initializing router...")

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		// 应用名称
		AppName: config.GetConfig().Server.Name,
		// 禁用启动消息
		DisableStartupMessage: false,
	})

	// 注册全局中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(compress.New())
	app.Use(requestid.New())
	app.Use(middleware.TraceID()) // TraceID中间件，用于链路追踪
	app.Use(recover.New())

	// Debug模式下启用日志中间件
	if config.GetConfig().Server.Debug {
		app.Use(fiberLogger.New(fiberLogger.Config{
			Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
		}))
	}

	// 注册路由
	health.InitialHealthRoutes(app)
	auth.InitialAuthRoutes(app)

	// 404处理
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code": 404,
			"msg":  "Route not found",
		})
	})

	// 启动服务
	port := config.GetConfig().Server.Port
	logger.Info("api", "Server starting on port %d...", port)

	if err := app.Listen(fmt.Sprintf("0.0.0.0:%d", port)); err != nil {
		logger.Error("fiber_listen", "Server failed to start: %v", err)
	}
}
