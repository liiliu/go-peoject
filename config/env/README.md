# 环境配置说明

## 📁 配置文件结构

```
config/
├── config.toml          # 默认配置（开发环境）
└── env/                 # 多环境配置目录
    ├── dev.toml         # 开发环境配置
    ├── test.toml        # 测试环境配置
    └── prod.toml        # 生产环境配置
```

## 🎯 环境说明

### dev.toml - 开发环境
- **用途**：本地开发使用
- **特点**：
  - 调试模式开启 `debug = true`
  - 启用自动数据库迁移 `autoMigrate = true`
  - 使用本地数据库和 Redis
  - 较小的连接池配置
  - 简单的 JWT 密钥（仅开发使用）

### test.toml - 测试环境
- **用途**：测试服务器部署
- **特点**：
  - 调试模式可选
  - 可启用自动迁移
  - 独立的数据库和 Redis DB
  - 中等连接池配置
  - 测试环境密钥

### prod.toml - 生产环境
- **用途**：正式生产环境
- **特点**：
  - **调试模式必须关闭** `debug = false`
  - **建议关闭自动迁移** `autoMigrate = false`
  - 使用生产数据库地址
  - Redis 必须设置密码
  - 最大连接池配置
  - **强密钥必须修改**

## 🚀 使用方法

### 本地开发

```bash
# 使用默认配置（config.toml）
go run main.go api

# 或者直接使用 dev.toml
cp config/env/dev.toml config/config.toml
go run main.go api
```

### 部署到测试环境

使用 `yun.sh` 部署时，配置环境变量：

```bash
# 在 yun.sh 中设置
CONFIG_ENV="test"

# 或运行时指定
CONFIG_ENV="test" bash yun.sh
```

脚本会自动复制 `config/env/test.toml` 为 `config/config.toml` 并打包部署。

### 部署到生产环境

```bash
# 在 yun.sh 中设置
CONFIG_ENV="prod"

# 或运行时指定
CONFIG_ENV="prod" bash yun.sh
```

## ⚙️ 配置项说明

### [server] 服务配置

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `application` | 应用标识 | "YOUR_PROJECT" |
| `name` | 项目名称 | "your_project" |
| `port` | 监听端口 | 8080 |
| `debug` | 调试模式 | true/false |
| `env` | 运行环境 | "DEV"/"TEST"/"PROD" |
| `logFileFolder` | 日志路径 | "./logs/log_error.log" |
| `tmpPath` | 临时目录 | "./tmp" |

### [mysqlMaster] 数据库配置

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `dsn` | 数据库连接串 | "user:pass@tcp(host:port)/db?..." |
| `maxIdleConn` | 最大空闲连接 | 10 |
| `maxOpenConn` | 最大打开连接 | 200 |
| `tablePrefix` | 表名前缀 | "app_" |
| `autoMigrate` | 自动迁移 | true/false |

### [redis] Redis 配置

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `host` | Redis 地址 | "127.0.0.1:6379" |
| `password` | Redis 密码 | "" |
| `db` | 数据库索引 | 0 |
| `poolSize` | 连接池大小 | 100 |
| `keyPrefix` | 缓存键前缀 | "APP_DEV" |

### [jwt] JWT 配置

| 配置项 | 说明 | 示例 |
|-------|------|------|
| `secret` | 加密密钥 | "your-secret-key" |
| `expireTime` | 过期时间（分钟） | 1440 |
| `maxRefreshTime` | 刷新时间（分钟） | 10080 |

## ⚠️ 重要提醒

### 开发环境（dev.toml）
- ✅ 可以使用简单的密码和密钥
- ✅ 可以开启 autoMigrate
- ✅ 可以使用本地服务

### 测试环境（test.toml）
- ⚠️ 使用独立的数据库和 Redis
- ⚠️ 密钥应该更复杂
- ✅ 可以开启 autoMigrate 便于测试

### 生产环境（prod.toml）
- ❌ **绝对不能**使用 `debug = true`
- ❌ **强烈建议**关闭 `autoMigrate = false`
- ❌ **务必修改** JWT secret 为强密钥
- ✅ **必须设置** Redis 密码
- ✅ **必须使用** 生产数据库地址
- ✅ **建议使用**环境变量或密钥管理服务

## 🔐 安全建议

### 1. 不要提交敏感信息到 Git

在 `.gitignore` 中添加：

```gitignore
# 生产环境配置不提交
config/env/prod.toml

# 或者所有环境配置都不提交（使用模板）
config/env/*.toml
!config/env/*.toml.example
```

### 2. 使用环境变量

生产环境建议使用环境变量覆盖配置：

```bash
export DB_PASSWORD="your-secure-password"
export JWT_SECRET="your-jwt-secret"
export REDIS_PASSWORD="your-redis-password"
```

### 3. 密钥强度

- 开发环境：可以简单
- 测试环境：中等复杂度
- 生产环境：**必须使用强密钥（32位以上随机字符）**

生成强密钥示例：
```bash
# 生成 32 字节随机密钥
openssl rand -base64 32
```

## 📚 相关文档

- [部署指南](../../docs/DEPLOYMENT.md)
- [项目 README](../../README.md)

---

**最后更新**：2025-10-16
