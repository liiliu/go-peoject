package main

import (
	"os"
	"your_project/commands/api"
	"your_project/commands/job"
	"your_project/commands/services"

	"github.com/urfave/cli/v2"
)

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
