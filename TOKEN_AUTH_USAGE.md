# Token 认证接口使用指南

## 概述

项目使用 JWT Token 进行身份认证。需要认证的接口必须在请求头中携带有效的 Token。

## 认证流程

### 1. 获取 Token

用户登录成功后会返回 Token：

```bash
POST http://localhost:8080/v1/auth/login
Content-Type: application/json

{
  "username": "your_username",
  "password": "your_password"
}
```

**响应示例：**
```json
{
  "result": 1,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "your_username",
    "user_id": 1
  }
}
```

### 2. 使用 Token 访问接口

在请求头中添加 Authorization 字段：

```bash
GET http://localhost:8080/v1/user/info
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**注意：** Token 前面必须加上 `Bearer ` 前缀（注意有空格）

## 已实现的需要 Token 的接口

### 用户信息相关

| 方法 | 路径 | 说明 | 需要 Token |
|------|------|------|-----------|
| GET | `/v1/user/info` | 获取当前用户信息 | ✅ |
| PUT | `/v1/user/update` | 更新用户信息 | ✅ |
| GET | `/v1/user/list` | 获取用户列表 | ✅ |

### 公开接口（不需要 Token）

| 方法 | 路径 | 说明 | 需要 Token |
|------|------|------|-----------|
| GET | `/v1/health` | 健康检查 | ❌ |
| POST | `/v1/auth/login` | 用户登录 | ❌ |
| POST | `/v1/auth/register` | 用户注册 | ❌ |
| GET | `/swagger/*` | Swagger 文档 | ❌ |

## 接口使用示例

### 1. 获取用户信息

```bash
curl -X GET http://localhost:8080/v1/user/info \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```

**响应：**
```json
{
  "result": 1,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": 1,
    "username": "zhangsan",
    "email": "zhangsan@example.com",
    "phone": "13800138000",
    "status": 1,
    "created_at": "2025-10-20T10:00:00Z",
    "updated_at": "2025-10-20T10:00:00Z"
  }
}
```

### 2. 更新用户信息

```bash
curl -X PUT http://localhost:8080/v1/user/update \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "new_email@example.com",
    "phone": "13900139000"
  }'
```

**响应：**
```json
{
  "result": 0,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 3. 获取用户列表（分页）

```bash
curl -X GET "http://localhost:8080/v1/user/list?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```

**响应：**
```json
{
  "result": 1,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "list": [
      {
        "id": 1,
        "username": "zhangsan",
        "email": "zhangsan@example.com",
        "phone": "13800138000",
        "status": 1
      }
    ],
    "total": 1,
    "page": 1
  }
}
```

## 在 Swagger 中测试需要 Token 的接口

### 步骤：

1. **获取 Token**
   - 访问 `http://localhost:8080/swagger/index.html`
   - 找到 `/v1/auth/login` 接口
   - 点击 "Try it out"
   - 输入用户名和密码
   - 点击 "Execute"
   - 复制返回的 token 值

2. **设置认证**
   - 点击页面右上角的 "Authorize" 🔓 按钮
   - 在弹出框中输入：`Bearer YOUR_TOKEN`（注意空格）
   - 点击 "Authorize" 按钮
   - 点击 "Close" 关闭

3. **测试接口**
   - 现在可以测试任何需要 Token 的接口
   - 展开接口 → "Try it out" → "Execute"

## 在代码中获取用户信息

在需要 Token 验证的接口中，可以从 Context 获取用户信息：

```go
func YourHandler(c *fiber.Ctx) error {
    // 获取用户ID
    userID, ok := c.Locals("user_id").(string)
    if !ok {
        return c.JSON(view.ErrorWithMsgCtx(c, view.CodeNoAuth, "身份鉴权失败"))
    }
    
    // 获取用户名
    username, ok := c.Locals("username").(string)
    if !ok {
        return c.JSON(view.ErrorWithMsgCtx(c, view.CodeNoAuth, "身份鉴权失败"))
    }
    
    logger.InfoWithTrace(c, "tag", "用户 %s (ID: %s) 访问接口", username, userID)
    
    // 业务逻辑...
}
```

## 错误响应

### Token 无效或过期

```json
{
  "result": 0,
  "code": 1007,
  "msg": "TOKEN失效",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 未提供 Token

```json
{
  "result": 0,
  "code": 1007,
  "msg": "TOKEN失效",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 身份鉴权失败

```json
{
  "result": 0,
  "code": 1006,
  "msg": "身份鉴权失败",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Token 验证中间件

项目使用 `CheckToken` 中间件进行 Token 验证：

```go
// 在路由中应用中间件
app.Use(middleware.CheckToken)

// 之后注册的所有路由都需要 Token 验证
user.InitialUserRoutes(app)
```

### 白名单配置

在 `app/api/middleware/auth.go` 中配置不需要验证的路由：

```go
var NoAuthUrls = []string{
    "/v1/auth/login",
    "/v1/auth/register",
    "/v1/health",
    "/swagger",  // Swagger 文档路径前缀
}
```

## 如何添加新的需要 Token 的接口

### 1. 创建接口函数

```go
// GetUserProfile 获取用户资料
// @Summary      获取用户资料
// @Description  获取指定用户的详细资料（需要Token）
// @Tags         用户
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  view.Result{data=model.User}  "成功"
// @Failure      1006 {object}  view.Result  "身份鉴权失败"
// @Failure      1007 {object}  view.Result  "TOKEN失效"
// @Router       /v1/user/profile/{id} [get]
func GetUserProfile(c *fiber.Ctx) error {
    // 获取当前登录用户信息
    userID := c.Locals("user_id").(string)
    username := c.Locals("username").(string)
    
    // 记录日志（带 traceId）
    logger.InfoWithTrace(c, "user", "用户 %s 查看资料", username)
    
    // 业务逻辑...
    
    // 返回响应（自动带 traceId）
    return c.JSON(view.SuccessWithCtx(c, data))
}
```

### 2. 注册路由

在 `CheckToken` 中间件**之后**注册路由：

```go
// router/router.go

// 应用 Token 验证中间件
app.Use(middleware.CheckToken)

// 注册需要 Token 验证的路由
user.InitialUserRoutes(app)
```

### 3. 重新生成 Swagger 文档

```bash
swag init
```

## 注意事项

1. **Token 格式**：必须使用 `Bearer ` 前缀，注意有空格
2. **Token 存储**：客户端应安全存储 Token（如 localStorage、sessionStorage）
3. **Token 过期**：Token 过期后需要重新登录获取新 Token
4. **安全性**：生产环境建议使用 HTTPS 传输
5. **用户信息**：Token 验证成功后，用户信息会自动存储到 Context
6. **日志追踪**：使用 `logger.InfoWithTrace()` 记录日志，自动包含 traceId

## 相关文件

- `app/api/middleware/auth.go` - Token 验证中间件
- `app/api/middleware/trace.go` - TraceID 中间件
- `app/api/v1/user/user.go` - 用户接口示例
- `library/jwt/jwt.go` - JWT Token 工具
- `router/router.go` - 路由配置

## 常见问题

### Q1: 如何让某个接口不需要 Token？
**A:** 在 `CheckToken` 中间件**之前**注册路由，或者将路径添加到 `NoAuthUrls` 白名单

### Q2: Token 存储在哪里？
**A:** 客户端负责存储，可以存在 localStorage、Cookie 或移动端的安全存储中

### Q3: 如何刷新 Token？
**A:** 可以实现一个 `/v1/auth/refresh` 接口用于刷新 Token

### Q4: 如何实现角色权限控制？
**A:** 可以扩展 Token 中的 Claims，添加角色信息，并创建新的权限验证中间件

## 完整流程示例

```bash
# 1. 用户登录获取 Token
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 2. 使用 Token 访问受保护的接口
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl -X GET http://localhost:8080/v1/user/info \
  -H "Authorization: Bearer $TOKEN"

# 3. 更新用户信息
curl -X PUT http://localhost:8080/v1/user/update \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"new@example.com"}'
```

## 总结

- ✅ 所有需要认证的接口都在 `CheckToken` 中间件之后注册
- ✅ Token 通过 `Authorization: Bearer {token}` 请求头传递
- ✅ 验证成功后自动将用户信息存入 Context
- ✅ 支持在 Swagger UI 中测试需要 Token 的接口
- ✅ 所有响应自动包含 traceId 用于链路追踪
