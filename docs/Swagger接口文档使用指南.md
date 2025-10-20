# Swagger 接口文档使用指南

## 概述

项目已集成 Swagger UI，提供可视化的 API 文档和在线测试功能。

## 访问 Swagger 文档

### 启动服务
```bash
go run main.go api
```

### 访问地址
服务启动后，在浏览器访问：
```
http://localhost:8080/swagger/index.html
```

## 主要功能

### 1. API 文档浏览
- ✅ 查看所有接口列表
- ✅ 查看接口详细信息（请求参数、响应格式、错误码等）
- ✅ 按 Tag 分组展示接口

### 2. 在线测试
- ✅ 直接在浏览器测试 API
- ✅ 填写参数后一键发送请求
- ✅ 查看实时响应结果

### 3. 模型查看
- ✅ 查看请求/响应数据结构
- ✅ 查看字段类型和说明

## 当前已集成接口

### 健康检查
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/v1/health` | 检查服务是否正常运行 |

### 认证接口
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/v1/auth/login` | 用户登录 |
| POST | `/v1/auth/register` | 用户注册 |

## 如何为新接口添加 Swagger 文档

### 1. 在接口函数上添加注释

```go
// Login 用户登录
// @Summary      用户登录
// @Description  用户通过用户名和密码登录，返回 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequest  true  "登录参数"
// @Success      200      {object}  view.Result{data=object{token=string,username=string,user_id=int}}  "登录成功"
// @Failure      1002     {object}  view.Result  "非法参数"
// @Failure      1004     {object}  view.Result  "参数验证失败"
// @Failure      1001     {object}  view.Result  "用户名或密码错误"
// @Router       /v1/auth/login [post]
func Login(c *fiber.Ctx) error {
    // 实现代码...
}
```

### 2. 注释字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| @Summary | 接口简短描述 | `用户登录` |
| @Description | 接口详细说明 | `用户通过用户名和密码登录，返回 JWT Token` |
| @Tags | 接口分组标签 | `认证` |
| @Accept | 接受的内容类型 | `json` |
| @Produce | 返回的内容类型 | `json` |
| @Param | 请求参数 | `request body request.LoginRequest true "登录参数"` |
| @Success | 成功响应 | `200 {object} view.Result "成功"` |
| @Failure | 失败响应 | `1001 {object} view.Result "失败"` |
| @Router | 路由路径和方法 | `/v1/auth/login [post]` |
| @Security | 安全认证 | `Bearer` |

### 3. 带认证的接口示例

对于需要 JWT Token 的接口：

```go
// GetUserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  view.Result{data=model.User}  "成功"
// @Failure      1006 {object}  view.Result  "身份鉴权失败"
// @Failure      1007 {object}  view.Result  "TOKEN失效"
// @Router       /v1/user/info [get]
func GetUserInfo(c *fiber.Ctx) error {
    // 实现代码...
}
```

### 4. 重新生成文档

每次修改注释后，需要重新生成文档：

```bash
swag init
```

或者使用开发时自动生成（需要安装 air）：

```bash
air
```

## Swagger 注释完整示例

```go
// CreateOrder 创建订单
// @Summary      创建订单
// @Description  用户创建新订单
// @Tags         订单
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request  body      request.CreateOrderRequest  true  "订单参数"
// @Success      200      {object}  view.Result{data=model.Order}  "创建成功"
// @Failure      1002     {object}  view.Result  "非法参数"
// @Failure      1004     {object}  view.Result  "参数验证失败"
// @Failure      9999     {object}  view.Result  "系统错误"
// @Router       /v1/order/create [post]
func CreateOrder(c *fiber.Ctx) error {
    // 实现代码...
}
```

## 高级配置

### 修改 Swagger 基础信息

在 `main.go` 文件顶部修改：

```go
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
```

### 配置生产环境地址

修改 `@host` 字段：
```go
// @host      api.example.com
```

## 使用技巧

### 1. 测试需要认证的接口

1. 先调用 `/v1/auth/login` 接口获取 Token
2. 点击右上角 "Authorize" 按钮
3. 输入 `Bearer {your_token}`（注意 Bearer 和 token 之间有空格）
4. 点击 "Authorize" 确认
5. 现在可以测试需要认证的接口了

### 2. 响应示例

所有接口都返回统一的响应格式：

```json
{
  "result": 1,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    // 实际数据
  }
}
```

### 3. TraceID 追踪

响应中包含 `trace_id` 字段，可用于日志追踪和问题排查。

## 常见问题

### Q1: 修改注释后没有生效？
**A:** 需要重新运行 `swag init` 生成文档

### Q2: 文档页面404？
**A:** 确保服务已启动，访问 `http://localhost:8080/swagger/index.html`

### Q3: 如何隐藏某些接口？
**A:** 不添加 Swagger 注释即可，或者在注释中使用 `@hidden`

### Q4: 如何测试上传文件的接口？
**A:** 使用 `@Param` 中的 `formData` 类型
```go
// @Param file formData file true "文件"
```

## 相关命令

```bash
# 安装 swag 命令行工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成/更新 Swagger 文档
swag init

# 查看 swag 版本
swag --version

# 格式化 Swagger 注释
swag fmt
```

## 参考资源

- [Swag GitHub](https://github.com/swaggo/swag)
- [Swagger 注释格式](https://github.com/swaggo/swag#declarative-comments-format)
- [Fiber Swagger](https://github.com/gofiber/swagger)
- [在线 Swagger 编辑器](https://editor.swagger.io/)

## 注意事项

1. **保持注释准确**：注释要与实际接口行为一致
2. **及时更新文档**：修改接口后记得更新注释
3. **合理分组**：使用 @Tags 合理组织接口
4. **错误码齐全**：列出所有可能的错误响应
5. **安全性**：生产环境可以考虑禁用 Swagger 或加上认证保护
