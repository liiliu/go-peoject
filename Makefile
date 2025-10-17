# Makefile for Go Project

# 项目名称
BINARY_NAME=your_project
MAIN_FILE=main.go

# 构建目录
BUILD_DIR=bin

# Go 命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# 运行命令
.PHONY: all build clean test run-api run-job run-worker run-grpc run-websocket dev fmt lint help

# 默认目标
all: build

# 构建项目
build:
	@echo "Building project..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 构建 Linux 版本（交叉编译）
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_FILE)
	@echo "Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux"

# 运行 API 服务
run-api:
	@echo "Starting API server..."
	$(GOCMD) run $(MAIN_FILE) api

# 运行定时任务
run-job:
	@echo "Starting job scheduler..."
	$(GOCMD) run $(MAIN_FILE) job

# 运行消息队列消费者
run-worker:
	@echo "Starting worker service..."
	$(GOCMD) run $(MAIN_FILE) worker

# 运行 gRPC 服务
run-grpc:
	@echo "Starting gRPC server..."
	$(GOCMD) run $(MAIN_FILE) grpc

# 运行 WebSocket 服务
run-websocket:
	@echo "Starting WebSocket server..."
	$(GOCMD) run $(MAIN_FILE) websocket

# 开发模式（需要安装 air: go install github.com/air-verse/air@latest）
dev:
	@echo "Starting development mode with hot reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

# 运行测试
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 运行测试并生成覆盖率报告
test-cover:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# 代码检查（需要安装 golangci-lint）
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# 清理编译文件
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# 帮助信息
help:
	@echo "Makefile Commands:"
	@echo "  make build       - Build the project"
	@echo "  make build-linux - Build for Linux (cross-compile)"
	@echo "  make run-api     - Run API server"
	@echo "  make run-job     - Run job scheduler"
	@echo "  make dev         - Start development mode with hot reload"
	@echo "  make test        - Run tests"
	@echo "  make test-cover  - Run tests with coverage report"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter"
	@echo "  make deps        - Download dependencies"
	@echo "  make clean       - Clean build files"
	@echo "  make help        - Show this help message"
