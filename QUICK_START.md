# 快速开始指南

## 🚀 5分钟快速部署

### 步骤1：克隆或移动项目

```bash
# 假设你已经将此模板移到了新的项目目录
cd your_new_project
```

### 步骤2：修改项目名称

**全局替换 `your_project` 为你的项目名：**

- 修改 `go.mod` 第一行的 module 名称
- 全局搜索替换代码中的 `your_project` import 路径

### 步骤3：配置数据库和Redis

编辑 `config/config.toml`：

```toml
[mysqlMaster]
dsn = "root:your_password@tcp(127.0.0.1:3306)/your_database?charset=utf8mb4&parseTime=True&loc=Local"

[redis]
host = "127.0.0.1:6379"
password = "your_redis_password"  # 如果有密码
```

### 步骤4：安装依赖

```bash
go mod tidy
go mod download
```

### 步骤5：创建数据库表

```sql
CREATE DATABASE IF NOT EXISTS your_database DEFAULT CHARACTER SET utf8mb4;

USE your_database;

-- 用户表示例
CREATE TABLE `app_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(100) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `status` int NOT NULL DEFAULT '1' COMMENT '状态 1-正常 0-禁用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 步骤6：运行项目

```bash
# 启动 API 服务
go run main.go api

# 或使用 Makefile
make run-api
```

### 步骤7：测试接口

```bash
# 健康检查
curl http://localhost:8080/v1/health

# 注册用户
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456",
    "email": "test@example.com"
  }'

# 登录获取Token
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

## 🐳 使用 Docker 快速部署

```bash
# 启动所有服务（包括 MySQL 和 Redis）
docker-compose up -d

# 查看日志
docker-compose logs -f api

# 停止服务
docker-compose down
```

## 📝 接下来做什么？

1. **添加新的API接口**
   - 在 `app/api/v1/` 创建新模块
   - 在 `router/router.go` 注册路由

2. **添加数据模型**
   - 在 `app/model/` 创建模型文件
   - 使用 GORM 进行数据库操作

3. **添加定时任务**
   - 在 `app/job/` 添加任务逻辑
   - 在 `app/job/job.go` 的 `scheduleJobs()` 注册cron表达式

4. **启用Token认证**
   - 在 `router/router.go` 中取消注释中间件
   - 或在特定路由组启用：`app.Use(middleware.CheckToken)`

5. **编写测试用例**
   - 参考 `app/api/v1/auth/login_test.go`
   - 运行测试：`go test ./...`

## 💡 常见问题

**Q: 如何热重载开发？**  
A: 安装 Air：`go install github.com/air-verse/air@latest`，然后运行 `make dev` 或 `air`

**Q: 如何修改端口？**  
A: 编辑 `config/config.toml` 中的 `port` 配置

**Q: 如何添加中间件？**  
A: 在 `app/api/middleware/` 创建中间件，在 `router/router.go` 使用 `app.Use()` 注册

**Q: 数据库迁移怎么做？**  
A: 可以使用 GORM 的 AutoMigrate 或第三方工具如 golang-migrate

## 📚 更多文档

查看完整的 README.md 了解更多功能和最佳实践。
