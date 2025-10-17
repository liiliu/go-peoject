# Go Starter Template - 文档中心

## 📚 文档导航

### 🚀 快速开始
- [项目快速开始](./快速开始.md) - 项目初始化、运行、部署
- [项目架构说明](./项目架构说明.md) - 目录结构、设计模式

### 📦 核心功能

#### 数据层
- [数据库操作指南](./数据库操作指南.md) - GORM 使用、事务、迁移
- [Redis缓存指南](./Redis缓存使用指南.md) - 缓存操作、分布式锁
- [MongoDB使用指南](./MongoDB使用指南.md) - MongoDB 操作

#### Web 开发
- [HTTP请求工具使用指南](./HTTP请求工具使用指南.md) - 发送 HTTP 请求
- [参数验证使用指南](./参数验证使用指南.md) - 请求参数验证
- [JWT认证使用指南](./JWT认证使用指南.md) - JWT 生成与验证

#### 文件处理
- [文件与数据处理工具使用指南](./文件与数据处理工具使用指南.md) - Excel、CSV、文件操作
- [Excel流式导出最佳实践](./Excel流式导出最佳实践.md) - 大数据量 Excel 导出

#### 工具类
- [工具函数使用指南](./工具函数使用指南.md) - 常用工具函数
- [日志使用指南](./日志使用指南.md) - 日志记录
- [验证码使用指南](./验证码使用指南.md) - 图形验证码

#### 消息队列
- [消息队列使用指南](./消息队列使用指南.md) - RabbitMQ、Kafka
- [异步任务队列](./异步任务队列.md) - 后台任务处理

#### 第三方服务
- [OSS对象存储指南](./OSS对象存储指南.md) - 阿里云OSS、MinIO
- [短信服务指南](./短信服务指南.md) - 短信发送
- [邮件服务指南](./邮件服务指南.md) - 邮件发送

### 🔧 部署运维
- [部署说明](./部署说明.md) - 生产环境部署
- [部署脚本说明](./部署脚本说明.md) - 自动化部署脚本

---

## 📖 类库速查

| 类库 | 说明 | 文档 |
|------|------|------|
| **database** | 数据库操作 (GORM) | [查看](./数据库操作指南.md) |
| **cache** | Redis 缓存 | [查看](./Redis缓存使用指南.md) |
| **http** | HTTP 客户端 | [查看](./HTTP请求工具使用指南.md) |
| **excel** | Excel 处理 | [查看](./文件与数据处理工具使用指南.md) |
| **csv** | CSV 处理 | [查看](./文件与数据处理工具使用指南.md) |
| **file** | 文件操作 | [查看](./文件与数据处理工具使用指南.md) |
| **jwt** | JWT 认证 | [查看](./JWT认证使用指南.md) |
| **validator** | 参数验证 | [查看](./参数验证使用指南.md) |
| **logger** | 日志记录 | [查看](./日志使用指南.md) |
| **util** | 工具函数 | [查看](./工具函数使用指南.md) |
| **captcha** | 验证码 | [查看](./验证码使用指南.md) |
| **queue** | 消息队列 | [查看](./消息队列使用指南.md) |
| **sms** | 短信服务 | [查看](./短信服务指南.md) |
| **email** | 邮件服务 | [查看](./邮件服务指南.md) |
| **oss** | 对象存储 | [查看](./OSS对象存储指南.md) |

---

## 🎯 常见场景

### 用户认证
```go
// JWT 生成
token, _ := jwt.GenerateToken(userID)

// JWT 验证
claims, _ := jwt.ParseToken(token)
```

### 数据库查询
```go
// 查询单条
var user model.User
database.DB.First(&user, 1)

// 分页查询
var users []model.User
database.DB.Limit(10).Offset(0).Find(&users)
```

### Redis 缓存
```go
// 设置缓存
cache.Set("key", "value", 3600)

// 获取缓存
value, _ := cache.Get("key")
```

### HTTP 请求
```go
// GET 请求
resp, _ := http.Get("https://api.example.com/users")

// POST 请求
resp, _ := http.Post("https://api.example.com/users", data)
```

### Excel 导出
```go
// 小数据量
headers := []string{"姓名", "年龄"}
data := [][]interface{}{{"张三", 25}}
excel.ExportData("users.xlsx", "用户", headers, data)

// 大数据量（流式）
excel.WriteFromCallback("users.xlsx", "用户", headers, func(page int) [][]interface{} {
    // 分页查询
    var users []User
    db.Limit(1000).Offset((page-1)*1000).Find(&users)
    // 转换数据...
})
```

---

## 🔥 推荐阅读顺序

### 新手入门
1. [快速开始](./快速开始.md) - 了解项目基本结构
2. [数据库操作指南](./数据库操作指南.md) - 学习数据库操作
3. [参数验证使用指南](./参数验证使用指南.md) - 学习参数验证
4. [JWT认证使用指南](./JWT认证使用指南.md) - 学习用户认证

### 进阶开发
1. [Redis缓存使用指南](./Redis缓存使用指南.md) - 使用缓存优化性能
2. [消息队列使用指南](./消息队列使用指南.md) - 处理异步任务
3. [Excel流式导出最佳实践](./Excel流式导出最佳实践.md) - 大数据导出
4. [HTTP请求工具使用指南](./HTTP请求工具使用指南.md) - 调用第三方 API

---

## 💡 最佳实践

### 项目结构
```
starter_template/
├── app/              # 应用层
│   ├── api/         # API 接口
│   ├── job/         # 定时任务
│   └── model/       # 数据模型
├── library/          # 工具库
├── config/           # 配置文件
├── docs/            # 文档
└── main.go          # 入口文件
```

### 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 Go 命名规范
- 添加必要的注释
- 错误处理不要忽略

### 性能优化
- 使用连接池
- 合理使用缓存
- 数据库索引优化
- 避免 N+1 查询

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

---

**更新日期**：2025-10-16  
**版本**：v1.0
