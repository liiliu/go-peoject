package main

import (
	"os"
	"your_project/commands/api"
	"your_project/commands/job"
	"your_project/commands/services"

	"github.com/urfave/cli/v2"
)

// @title           Your Project API
// @version         1.0
// @description     基于 Fiber 的 Go Web 应用 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 输入 Bearer token，格式：Bearer {token}

func main() {
	cmd := cli.NewApp()
	cmd.Name = "your_project"
	cmd.Usage = "Your Project - A Go web application based on Fiber"
	cmd.Version = "1.0.0"

	cmd.Commands = []*cli.Command{
		{
			Name:  "api",
			Usage: "Run API server",
			Action: func(c *cli.Context) error {
				api.Run()
				select {} // Keep the server running
			},
		},
		{
			Name:  "job",
			Usage: "Run scheduled job server",
			Action: func(c *cli.Context) error {
				job.Run()
				select {} // Keep the job scheduler running
			},
		},
		{
			Name:  "websocket",
			Usage: "Run WebSocket server",
			Action: func(c *cli.Context) error {
				services.RunWebSocket()
				return nil
			},
		},
	}

	if err := cmd.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
