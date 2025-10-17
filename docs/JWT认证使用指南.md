# JWT 认证使用指南

## 📋 概述

项目使用 **JWT (JSON Web Token)** 实现用户认证，提供无状态的分布式认证方案。

---

## 🚀 快速开始

### 1. 生成 Token

```go
import "your_project/library/jwt"

// 生成 Token
token, err := jwt.GenerateToken(userID)
if err != nil {
    return err
}

logger.Info("jwt", "Token: %s", token)
```

### 2. 解析 Token

```go
// 解析 Token
claims, err := jwt.ParseToken(token)
if err != nil {
    return errors.New("Token 无效")
}

userID := claims.UserID
logger.Info("jwt", "用户ID: %d", userID)
```

---

## 📝 完整认证流程

### 1. 用户登录

```go
// app/api/v1/auth/login.go
package auth

import (
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"
    "your_project/app/model"
    "your_project/library/database"
    "your_project/library/jwt"
    "your_project/library/validator"
)

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token string      `json:"token"`
    User  *model.User `json:"user"`
}

func Login(c *fiber.Ctx) error {
    var req LoginRequest
    
    // 解析参数
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "code": 400,
            "msg":  "参数错误",
        })
    }
    
    // 验证参数
    if err := validator.Validate.Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "code": 400,
            "msg":  validator.TranslateError(err),
        })
    }
    
    // 查询用户
    var user model.User
    if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
        return c.Status(401).JSON(fiber.Map{
            "code": 401,
            "msg":  "用户名或密码错误",
        })
    }
    
    // 验证密码
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return c.Status(401).JSON(fiber.Map{
            "code": 401,
            "msg":  "用户名或密码错误",
        })
    }
    
    // 生成 Token
    token, err := jwt.GenerateToken(user.ID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "code": 500,
            "msg":  "生成 Token 失败",
        })
    }
    
    return c.JSON(fiber.Map{
        "code": 200,
        "msg":  "登录成功",
        "data": LoginResponse{
            Token: token,
            User:  &user,
        },
    })
}
```

### 2. 认证中间件

```go
// app/api/middleware/auth.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "strings"
    "your_project/library/jwt"
)

// Auth JWT 认证中间件
func Auth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 获取 Token
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "code": 401,
                "msg":  "未登录",
            })
        }
        
        // 格式：Bearer <token>
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            return c.Status(401).JSON(fiber.Map{
                "code": 401,
                "msg":  "Token 格式错误",
            })
        }
        
        token := parts[1]
        
        // 解析 Token
        claims, err := jwt.ParseToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "code": 401,
                "msg":  "Token 无效或已过期",
            })
        }
        
        // 将用户ID保存到上下文
        c.Locals("user_id", claims.UserID)
        
        return c.Next()
    }
}
```

### 3. 使用认证中间件

```go
// commands/api/api.go
v1 := app.Group("/v1")

// 公开接口（不需要认证）
v1.Post("/auth/login", auth.Login)
v1.Post("/auth/register", auth.Register)

// 需要认证的接口
userGroup := v1.Group("/user")
userGroup.Use(middleware.Auth())  // 使用认证中间件

userGroup.Get("/profile", user.GetProfile)
userGroup.Put("/profile", user.UpdateProfile)
userGroup.Get("/orders", order.GetUserOrders)
```

### 4. 在控制器中获取用户ID

```go
// app/api/v1/user/user.go
func GetProfile(c *fiber.Ctx) error {
    // 从上下文获取用户ID
    userID := c.Locals("user_id").(uint)
    
    // 查询用户信息
    var user model.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{
            "code": 404,
            "msg":  "用户不存在",
        })
    }
    
    return c.JSON(fiber.Map{
        "code": 200,
        "data": user,
    })
}
```

---

## 🎯 进阶功能

### 1. Token 刷新

```go
// 刷新 Token
func RefreshToken(c *fiber.Ctx) error {
    // 从当前 Token 获取用户ID
    userID := c.Locals("user_id").(uint)
    
    // 生成新 Token
    newToken, err := jwt.GenerateToken(userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "code": 500,
            "msg":  "刷新失败",
        })
    }
    
    return c.JSON(fiber.Map{
        "code": 200,
        "data": fiber.Map{
            "token": newToken,
        },
    })
}
```

### 2. Token 黑名单（登出）

```go
// 登出（将 Token 加入黑名单）
func Logout(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    parts := strings.SplitN(authHeader, " ", 2)
    token := parts[1]
    
    // 解析 Token 获取过期时间
    claims, _ := jwt.ParseToken(token)
    expiration := int(claims.ExpiresAt.Time.Unix() - time.Now().Unix())
    
    // 将 Token 加入黑名单（Redis）
    key := fmt.Sprintf("blacklist:%s", token)
    cache.Set(key, "1", expiration)
    
    return c.JSON(fiber.Map{
        "code": 200,
        "msg":  "登出成功",
    })
}

// 认证中间件中检查黑名单
func Auth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // ... 获取 Token ...
        
        // 检查是否在黑名单中
        key := fmt.Sprintf("blacklist:%s", token)
        exists, _ := cache.Exists(key)
        if exists {
            return c.Status(401).JSON(fiber.Map{
                "code": 401,
                "msg":  "Token 已失效",
            })
        }
        
        // ... 解析 Token ...
    }
}
```

### 3. 自定义 Claims

```go
// 扩展 Claims
type CustomClaims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// 生成带自定义信息的 Token
func GenerateCustomToken(user *model.User) (string, error) {
    claims := CustomClaims{
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key"))
}
```

### 4. 权限验证中间件

```go
// 权限验证
func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(uint)
        
        // 查询用户角色
        var user model.User
        database.DB.Select("role").First(&user, userID)
        
        if user.Role != role {
            return c.Status(403).JSON(fiber.Map{
                "code": 403,
                "msg":  "权限不足",
            })
        }
        
        return c.Next()
    }
}

// 使用
adminGroup := v1.Group("/admin")
adminGroup.Use(middleware.Auth())
adminGroup.Use(middleware.RequireRole("admin"))
```

---

## 💡 实战示例

### 示例 1：完整的认证系统

```go
// 注册
func Register(c *fiber.Ctx) error {
    var req RegisterRequest
    
    // 解析和验证参数
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"code": 400, "msg": "参数错误"})
    }
    
    if err := validator.Validate.Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{"code": 400, "msg": validator.TranslateError(err)})
    }
    
    // 检查用户名是否存在
    var count int64
    database.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
    if count > 0 {
        return c.Status(400).JSON(fiber.Map{"code": 400, "msg": "用户名已存在"})
    }
    
    // 密码加密
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    // 创建用户
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: string(hashedPassword),
        Status:   1,
    }
    
    if err := database.DB.Create(user).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"code": 500, "msg": "注册失败"})
    }
    
    // 生成 Token
    token, _ := jwt.GenerateToken(user.ID)
    
    return c.JSON(fiber.Map{
        "code": 200,
        "msg":  "注册成功",
        "data": fiber.Map{
            "token": token,
            "user":  user,
        },
    })
}
```

### 示例 2：找回密码

```go
// 发送重置密码邮件
func SendResetPasswordEmail(c *fiber.Ctx) error {
    email := c.FormValue("email")
    
    // 查询用户
    var user model.User
    if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return c.JSON(fiber.Map{"code": 404, "msg": "邮箱不存在"})
    }
    
    // 生成重置Token（30分钟有效）
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "type":    "reset_password",
        "exp":     time.Now().Add(30 * time.Minute).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    resetToken, _ := token.SignedString([]byte("your-secret-key"))
    
    // 发送邮件（包含重置链接）
    resetURL := fmt.Sprintf("https://example.com/reset-password?token=%s", resetToken)
    // email.Send(user.Email, "重置密码", resetURL)
    
    return c.JSON(fiber.Map{"code": 200, "msg": "重置邮件已发送"})
}

// 重置密码
func ResetPassword(c *fiber.Ctx) error {
    var req struct {
        Token    string `json:"token"`
        Password string `json:"password"`
    }
    
    c.BodyParser(&req)
    
    // 解析 Token
    claims, err := jwt.ParseToken(req.Token)
    if err != nil {
        return c.JSON(fiber.Map{"code": 400, "msg": "Token 无效或已过期"})
    }
    
    // 检查 Token 类型
    if claims["type"] != "reset_password" {
        return c.JSON(fiber.Map{"code": 400, "msg": "Token 类型错误"})
    }
    
    // 更新密码
    userID := uint(claims["user_id"].(float64))
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    
    database.DB.Model(&model.User{}).Where("id = ?", userID).Update("password", string(hashedPassword))
    
    return c.JSON(fiber.Map{"code": 200, "msg": "密码重置成功"})
}
```

### 示例 3：多端登录

```go
// 支持Web、APP、小程序等多端登录
func Login(c *fiber.Ctx) error {
    var req LoginRequest
    c.BodyParser(&req)
    
    // ... 验证用户名密码 ...
    
    // 获取设备类型
    deviceType := c.Get("Device-Type", "web")  // web, app, miniapp
    
    // 生成 Token（不同设备不同过期时间）
    var expiration time.Duration
    switch deviceType {
    case "web":
        expiration = 24 * time.Hour       // Web 1天
    case "app":
        expiration = 30 * 24 * time.Hour  // APP 30天
    case "miniapp":
        expiration = 7 * 24 * time.Hour   // 小程序 7天
    default:
        expiration = 24 * time.Hour
    }
    
    // 生成 Token
    token, _ := jwt.GenerateTokenWithExpiration(user.ID, expiration)
    
    // 保存 Token 到 Redis（用于单点登录控制）
    key := fmt.Sprintf("token:%d:%s", user.ID, deviceType)
    cache.Set(key, token, int(expiration.Seconds()))
    
    return c.JSON(fiber.Map{
        "code": 200,
        "data": fiber.Map{
            "token": token,
        },
    })
}
```

---

## ⚠️ 安全建议

### 1. Secret Key 安全

```go
// ❌ 不要硬编码
const secretKey = "my-secret-key"

// ✅ 从配置文件读取
secretKey := config.GetConfig().JWT.Secret
```

### 2. Token 过期时间

```go
// 根据业务场景设置合理的过期时间
- Web端：24小时
- APP端：7-30天
- 支付等敏感操作：15-30分钟
```

### 3. HTTPS

```go
// 生产环境必须使用 HTTPS
// Token 通过 HTTP 传输容易被截获
```

### 4. 刷新 Token 机制

```go
// 实现双 Token 机制
type TokenPair struct {
    AccessToken  string `json:"access_token"`   // 短期（1小时）
    RefreshToken string `json:"refresh_token"`  // 长期（30天）
}

// AccessToken 过期时，使用 RefreshToken 获取新的 AccessToken
```

---

## 📚 相关资源

- **JWT 官网**：https://jwt.io/
- **golang-jwt**：https://github.com/golang-jwt/jwt

---

**更新日期**：2025-10-16  
**版本**：v1.0
