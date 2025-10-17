# Go 项目脚手架模板

基于 Fiber + GORM + Redis 的 Go Web 项目脚手架，包含完整的项目结构和常用组件封装。

## 项目特性

- ✅ **Web 框架**：基于 Fiber v2（高性能 HTTP 框架）
- ✅ **数据库 ORM**：GORM（支持 MySQL/PostgreSQL）
- ✅ **数据库迁移**：集成 GORM AutoMigrate，支持自动创建/更新表结构
- ✅ **缓存**：Redis 封装
- ✅ **认证**：JWT Token 认证中间件
- ✅ **日志**：基于 Zap 的结构化日志
- ✅ **配置管理**：TOML 配置文件
- ✅ **定时任务**：基于 cron 的任务调度
- ✅ **测试**：单元测试示例
- ✅ **工具函数**：常用工具方法封装

## 目录结构

```
.
├── app/                    # 应用层
│   ├── api/               # API 处理器
│   │   ├── middleware/    # 中间件
│   │   └── v1/           # API v1 版本
│   │       ├── auth/     # 认证相关接口
│   │       └── health/   # 健康检查
│   ├── job/              # 定时任务
│   ├── model/            # 数据模型
│   ├── request/          # 请求参数结构体
│   └── view/             # 响应结构体
├── commands/              # 命令入口
│   ├── api/              # API 服务启动
│   └── job/              # 定时任务启动
├── config/                # 配置文件
│   ├── config.toml       # 主配置文件
│   └── env/              # 多环境配置（可选）
├── library/               # 核心组件库
│   ├── cache/            # Redis 缓存封装
│   ├── common/           # 公共常量
│   ├── config/           # 配置加载
│   ├── database/         # 数据库连接
│   ├── jwt/              # JWT 认证
│   ├── logger/           # 日志组件
│   └── util/             # 工具函数
├── router/                # 路由注册
├── tmp/                   # 临时文件目录
├── logs/                  # 日志文件目录
├── .gitignore            # Git 忽略文件
├── go.mod                # Go 模块依赖
├── main.go               # 程序入口
├── Makefile              # 构建脚本
└── README.md             # 项目说明文档
```

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置文件

编辑 `config/config.toml`，修改数据库、Redis 等配置：

```toml
[server]
name = "your_project"
port = 8080
debug = true
env = "DEV"
logFileFolder = "./logs/log_error.log"

[mysqlMaster]
dsn = "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
maxIdleConn = 10
maxOpenConn = 200
tablePrefix = "app_"
autoMigrate = true  # 开发环境启用自动迁移

[redis]
host = "127.0.0.1:6379"
password = ""
db = 0
keyPrefix = "APP_DEV"
```

### 3. 运行项目

**启动 API 服务：**
```bash
go run main.go api
```

**启动定时任务：**
```bash
go run main.go job
```

**启动其他服务：**
```bash
go run main.go worker     # 启动消息队列消费者
go run main.go grpc       # 启动 gRPC 服务
go run main.go websocket  # 启动 WebSocket 服务
```

**使用 Makefile：**
```bash
make run-api        # 启动 API 服务
make run-job        # 启动定时任务
make run-worker     # 启动消息队列消费者
make run-grpc       # 启动 gRPC 服务
make run-websocket  # 启动 WebSocket 服务
make build          # 编译项目
make test           # 运行测试
```

**查看所有可用命令：**
```bash
go run main.go --help
```

### 4. 测试接口

```bash
# 健康检查
curl http://localhost:8080/v1/health

# 获取 Token（示例）
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

## 部署到服务器

### 快速部署

```bash
# 部署到测试环境
bash yun-test.sh

# 部署到生产环境（需确认）
bash yun-prod.sh
```

### 部署脚本说明

项目提供了环境特定的部署脚本：

| 脚本 | 环境 | 用途 |
|------|------|------|
| `yun-test.sh` | 测试环境 | 编译、打包并部署到测试服务器 |
| `yun-prod.sh` | 生产环境 | 编译、打包并部署到生产服务器 |
| `restart.sh` | 服务器端 | 服务管理（启动/停止/监控） |

### 配置部署脚本

**测试环境** - 编辑 `yun-test.sh`：

```bash
SERVER_HOST="119.91.131.93"      # 测试服务器 IP
SERVER_PORT="61022"               # SSH 端口
SERVER_USER="root"                # 登录用户
SERVER_PATH="/data/server/..."    # 部署目录
```

**生产环境** - 编辑 `yun-prod.sh`：

```bash
SERVER_HOST="your-prod-ip"        # 生产服务器 IP（必须修改）
SERVER_PORT=""                    # SSH 端口
SERVER_USER="root"
SERVER_PATH="/data/server/..."    # 部署目录（必须修改）
```

**环境配置文件**自动选择：
- `yun-test.sh` → 使用 `config/env/test.toml`
- `yun-prod.sh` → 使用 `config/env/prod.toml`

### 服务器端管理

```bash
# 启动服务
bash restart.sh start

# 停止服务
bash restart.sh stop

# 查看状态
bash restart.sh status

# 查看日志
bash restart.sh logs api 100

# 监控模式（自动重启异常服务）
bash restart.sh watch
```

**详细文档**：
- [完整部署指南](docs/DEPLOYMENT.md) - 详细的部署流程和配置
- [部署脚本手册](docs/DEPLOYMENT_SCRIPTS.md) - 脚本使用说明和故障排查

## 核心组件使用

### 配置管理

```go
import "your_project/library/config"

// 获取配置
cfg := config.GetConfig()
port := cfg.Server.Port
```

### 数据库操作

```go
import "your_project/library/database"

// 获取数据库实例
db := database.NewEngine()

// GORM 操作
var user model.User
db.Where("id = ?", 1).First(&user)
```

### 数据库迁移（AutoMigrate）

**注册模型（在 `commands/api/api.go` 或 `commands/job/job.go` 中）：**
```go
import (
    "your_project/app/model"
    "your_project/library/database"
)

func registerModels() {
    database.RegisterModels(
        &model.User{},
        &model.Article{},
        // 添加更多模型...
    )
}
```

**配置启用（在 `config/config.toml` 中）：**
```toml
[mysqlMaster]
autoMigrate = true  # 开发环境启用，生产环境建议关闭
```

**详细文档**：查看 [`docs/AUTO_MIGRATE.md`](docs/AUTO_MIGRATE.md) 了解完整使用指南

### Redis 缓存

```go
import "your_project/library/cache"

// 保存数据
cache.Save("key", data, 5*time.Minute)

// 获取数据
var result YourStruct
cache.Get("key", &result)

// 保存字符串
cache.SaveString("key", "value", time.Hour)

// 获取字符串
value := cache.GetString("key")
```

### JWT 认证

```go
import "your_project/library/jwt"

// 生成 Token
j := jwt.NewJWT()
token := j.IssueToken(userID, isSuper)

// 解析 Token
claims, err := j.ParserToken(c)
if err != nil {
    return err
}
userID := claims.UserID
```

### 日志记录

```go
import "your_project/library/logger"

// 普通日志
logger.Info("用户登录成功", userID)
logger.Debug("调试信息: %v", data)

// 错误日志（带标签）
logger.Error("API错误", "请求失败: %v", err)
```

### 定时任务

在 `app/job/job.go` 中添加定时任务：

```go
func scheduleJobs() {
    // 每分钟执行
    cronEngine.AddFunc("0 * * * * *", func() {
        YourTaskFunction()
    })
    
    // 每天凌晨3点执行
    cronEngine.AddFunc("0 0 3 * * *", func() {
        DailyTask()
    })
}
```

## 开发指南

### 添加新接口

1. 在 `app/api/v1/` 下创建新模块
2. 定义请求/响应结构体（`app/request/`, `app/view/`）
3. 实现业务逻辑
4. 在 `router/router.go` 中注册路由

### 添加数据模型

在 `app/model/` 下创建模型文件：

```go
package model

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex;size:50"`
    Password  string    `gorm:"size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (User) TableName() string {
    return "users"
}
```

### 编写测试

在对应文件旁创建 `*_test.go` 文件：

```go
package auth

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
    // 测试逻辑
    assert.Equal(t, expected, actual)
}
```

## 常用命令

```bash
# 开发
make dev          # 热重载开发模式（需安装 air）
make run-api      # 运行 API 服务
make run-job      # 运行定时任务

# 构建
make build        # 编译项目
make build-linux  # 交叉编译 Linux 版本

# 测试
make test         # 运行所有测试
make test-cover   # 测试覆盖率

# 代码质量
make fmt          # 格式化代码
make lint         # 代码检查

# 清理
make clean        # 清理编译文件
```

## 部署

### Docker 部署（可选）

```bash
docker build -t your_project .
docker run -p 8080:8080 your_project
```

### 直接部署

```bash
make build
./bin/your_project api
```

## 常见问题

**Q: 如何修改项目名称？**  
A: 全局替换 `your_project` 为你的项目名，修改 `go.mod` 中的 module 名称。

**Q: 如何添加新的中间件？**  
A: 在 `app/api/middleware/` 下创建中间件，在 `router/router.go` 中使用 `app.Use()` 注册。

**Q: 日志文件在哪里？**  
A: 日志文件保存在 `logs/` 目录下，可在配置文件中修改路径。

## License

MIT License
