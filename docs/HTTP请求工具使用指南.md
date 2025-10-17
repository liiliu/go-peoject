# HTTP 请求工具使用指南

## 📋 概述

项目使用 **`go-resty/resty`** 作为 HTTP 客户端库，并在 `library/http` 中进行了二次封装，提供更简洁易用的 API。

---

## 🚀 快速开始

### 安装依赖

```bash
go get github.com/go-resty/resty/v2
```

### 简单使用

```go
import "your_project/library/http"

// GET 请求
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    logger.Error("http", "Request failed: %v", err)
    return
}

logger.Info("http", "Response: %s", resp.String())
```

---

## 📚 基础用法

### 1. GET 请求

#### 简单 GET

```go
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    return err
}

// 获取响应内容
body := resp.String()
statusCode := resp.StatusCode()
```

#### 带参数的 GET

```go
params := map[string]string{
    "page":   "1",
    "size":   "10",
    "status": "active",
}

resp, err := http.GetWithParams("https://api.example.com/users", params)
// 实际请求: https://api.example.com/users?page=1&size=10&status=active
```

#### 带自定义 Header

```go
headers := map[string]string{
    "Authorization": "Bearer your-token",
    "Custom-Header": "custom-value",
}

resp, err := http.GetWithHeaders("https://api.example.com/users", headers)
```

---

### 2. POST 请求

#### POST JSON

```go
// 方式 1：使用 map
data := map[string]interface{}{
    "name":  "张三",
    "email": "zhangsan@example.com",
    "age":   25,
}
resp, err := http.Post("https://api.example.com/users", data)

// 方式 2：使用结构体
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}
user := CreateUserRequest{
    Name:  "李四",
    Email: "lisi@example.com",
    Age:   30,
}
resp, err := http.Post("https://api.example.com/users", user)
```

#### POST JSON 并解析响应

```go
// 请求数据
reqData := map[string]interface{}{
    "username": "admin",
    "password": "123456",
}

// 响应结构体
type LoginResponse struct {
    Token  string `json:"token"`
    UserID int    `json:"user_id"`
}

var result LoginResponse
resp, err := http.PostJSON("https://api.example.com/login", reqData, &result)

if err != nil {
    return err
}

logger.Info("http", "Login success, token: %s", result.Token)
```

#### POST 表单

```go
formData := map[string]string{
    "username": "admin",
    "password": "123456",
    "remember": "1",
}

resp, err := http.PostForm("https://api.example.com/login", formData)
```

---

### 3. 其他 HTTP 方法

#### PUT 请求

```go
updateData := map[string]interface{}{
    "name":   "新名字",
    "status": "active",
}

resp, err := http.Put("https://api.example.com/users/1", updateData)
```

#### DELETE 请求

```go
resp, err := http.Delete("https://api.example.com/users/1")
```

#### PATCH 请求

```go
patchData := map[string]interface{}{
    "status": "inactive",
}

resp, err := http.Patch("https://api.example.com/users/1", patchData)
```

---

## 🎯 高级用法

### 1. 链式调用（推荐）

```go
resp, err := http.Client.R().
    SetHeader("Authorization", "Bearer token").
    SetQueryParam("status", "active").
    SetQueryParam("page", "1").
    SetBody(map[string]string{
        "name": "test",
    }).
    Post("https://api.example.com/users")
```

### 2. 使用请求构建器

```go
req := http.NewRequest().
    SetAuthToken("your-token").
    SetQueryParam("page", "1").
    SetQueryParam("size", "10").
    SetBody(data)

resp, err := req.Execute("POST", "https://api.example.com/orders")

// 或者
resp, err := req.Send()
```

### 3. 设置 Bearer Token

```go
resp, err := http.Client.R().
    SetAuthToken("your-jwt-token").
    Get("https://api.example.com/protected")
```

### 4. 基础认证

```go
resp, err := http.Client.R().
    SetBasicAuth("username", "password").
    Get("https://api.example.com/admin")
```

### 5. 自定义超时

**注意**：resty 的超时需要在 Client 级别设置，不能在单个请求上设置。

```go
// 方式 1：创建新客户端并设置超时
client := http.NewClient()
client.SetTimeout(5 * time.Second)
resp, err := client.R().Get("https://api.example.com/slow-endpoint")

// 方式 2：使用配置创建客户端
client := http.NewClientWithConfig(http.Config{
    Timeout: 5 * time.Second,
})
resp, err := client.R().Get("https://api.example.com/slow-endpoint")

// 方式 3：修改全局客户端（影响所有后续请求）
http.Client.SetTimeout(5 * time.Second)
resp, err := http.Get("https://api.example.com/slow-endpoint")
```

### 6. 错误处理

```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

var result SuccessResponse
var errResp ErrorResponse

resp, err := http.Client.R().
    SetResult(&result).      // 成功时解析到这里
    SetError(&errResp).      // 失败时解析到这里
    Post("https://api.example.com/users", data)

if err != nil {
    logger.Error("http", "Request error: %v", err)
    return
}

if !http.IsSuccess(resp) {
    logger.Error("http", "API error: %s", errResp.Message)
    return
}

logger.Info("http", "Success: %+v", result)
```

---

## 📁 文件操作

### 1. 文件上传

```go
// 单文件上传
resp, err := http.Client.R().
    SetFile("file", "/path/to/file.jpg").
    SetFormData(map[string]string{
        "title":       "My Photo",
        "description": "Beautiful sunset",
    }).
    Post("https://api.example.com/upload")

// 多文件上传
resp, err := http.Client.R().
    SetFiles(map[string]string{
        "file1": "/path/to/file1.jpg",
        "file2": "/path/to/file2.jpg",
    }).
    Post("https://api.example.com/upload/multiple")
```

### 2. 文件下载

```go
// 下载到文件
resp, err := http.Client.R().
    SetOutput("/path/to/save/file.pdf").
    Get("https://example.com/files/document.pdf")

if err != nil {
    logger.Error("http", "Download failed: %v", err)
    return
}

logger.Info("http", "File downloaded successfully")
```

---

## 🔧 自定义配置

### 创建自定义客户端

```go
import (
    "time"
    "your_project/library/http"
)

// 创建自定义配置的客户端
client := http.NewClientWithConfig(http.Config{
    Timeout:       10 * time.Second,
    RetryCount:    5,
    RetryWaitTime: 2 * time.Second,
    BaseURL:       "https://api.example.com",
    Headers: map[string]string{
        "Authorization": "Bearer token",
        "Custom-Header": "value",
    },
    Debug: true,  // 开启调试模式
})

// 使用自定义客户端
resp, err := client.R().Get("/users")
```

### 全局客户端配置

```go
// 修改全局客户端配置
http.Client.SetTimeout(30 * time.Second)
http.Client.SetRetryCount(3)
http.Client.SetBaseURL("https://api.example.com")

// 设置全局 Header
http.Client.SetHeader("Authorization", "Bearer token")
http.Client.SetHeaders(map[string]string{
    "Custom-Header-1": "value1",
    "Custom-Header-2": "value2",
})
```

---

## 💡 实战示例

### 示例 1：调用第三方 API

```go
package service

import (
    "your_project/library/http"
    "your_project/library/logger"
)

type WeatherResponse struct {
    City        string  `json:"city"`
    Temperature float64 `json:"temperature"`
    Description string  `json:"description"`
}

func GetWeather(city string) (*WeatherResponse, error) {
    var result WeatherResponse
    
    resp, err := http.Client.R().
        SetQueryParam("city", city).
        SetQueryParam("key", "your-api-key").
        SetResult(&result).
        Get("https://api.weather.com/v1/current")
    
    if err != nil {
        logger.Error("weather", "Failed to get weather: %v", err)
        return nil, err
    }
    
    if !http.IsSuccess(resp) {
        logger.Error("weather", "API error: %d", resp.StatusCode())
        return nil, fmt.Errorf("API error: %d", resp.StatusCode())
    }
    
    return &result, nil
}
```

### 示例 2：微信小程序登录

```go
type WechatLoginResponse struct {
    OpenID     string `json:"openid"`
    SessionKey string `json:"session_key"`
    UnionID    string `json:"unionid"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}

func WechatLogin(code string) (*WechatLoginResponse, error) {
    appID := config.GetConfig().Wechat.AppID
    secret := config.GetConfig().Wechat.Secret
    
    var result WechatLoginResponse
    
    resp, err := http.Client.R().
        SetQueryParams(map[string]string{
            "appid":      appID,
            "secret":     secret,
            "js_code":    code,
            "grant_type": "authorization_code",
        }).
        SetResult(&result).
        Get("https://api.weixin.qq.com/sns/jscode2session")
    
    if err != nil {
        return nil, err
    }
    
    if result.ErrCode != 0 {
        return nil, fmt.Errorf("wechat error: %s", result.ErrMsg)
    }
    
    return &result, nil
}
```

### 示例 3：支付宝支付回调

```go
func AlipayCallback(c *fiber.Ctx) error {
    // 获取表单数据
    formData := make(map[string]string)
    c.Request().PostArgs().VisitAll(func(key, value []byte) {
        formData[string(key)] = string(value)
    })
    
    // 验证签名（调用支付宝 API）
    resp, err := http.PostForm(
        "https://openapi.alipay.com/gateway.do",
        formData,
    )
    
    if err != nil {
        logger.Error("alipay", "Callback verification failed: %v", err)
        return c.SendString("fail")
    }
    
    // 处理业务逻辑...
    return c.SendString("success")
}
```

### 示例 4：短信发送（阿里云）

```go
type AliyunSmsResponse struct {
    Code      string `json:"Code"`
    Message   string `json:"Message"`
    RequestId string `json:"RequestId"`
}

func SendSms(phone, code string) error {
    params := map[string]string{
        "PhoneNumbers":  phone,
        "SignName":      "您的签名",
        "TemplateCode":  "SMS_123456",
        "TemplateParam": fmt.Sprintf(`{"code":"%s"}`, code),
    }
    
    var result AliyunSmsResponse
    
    resp, err := http.Client.R().
        SetHeader("Authorization", "Bearer "+getAliyunToken()).
        SetQueryParams(params).
        SetResult(&result).
        Post("https://dysmsapi.aliyuncs.com/")
    
    if err != nil {
        return err
    }
    
    if result.Code != "OK" {
        return fmt.Errorf("sms error: %s", result.Message)
    }
    
    return nil
}
```

### 示例 5：并发请求

```go
import "sync"

func FetchMultipleUsers(userIDs []int) ([]User, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    users := make([]User, 0, len(userIDs))
    
    for _, id := range userIDs {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            var user User
            _, err := http.Client.R().
                SetResult(&user).
                Get(fmt.Sprintf("https://api.example.com/users/%d", userID))
            
            if err == nil {
                mu.Lock()
                users = append(users, user)
                mu.Unlock()
            }
        }(id)
    }
    
    wg.Wait()
    return users, nil
}
```

---

## ⚠️ 注意事项

### 1. 超时设置

生产环境建议设置合理的超时时间：

```go
http.Client.SetTimeout(10 * time.Second)  // 全局设置

// 或单个请求设置
resp, err := http.Client.R().
    SetTimeout(5 * time.Second).
    Get(url)
```

### 2. 重试机制

```go
// 设置重试次数和等待时间
http.Client.SetRetryCount(3)
http.Client.SetRetryWaitTime(1 * time.Second)
http.Client.SetRetryMaxWaitTime(5 * time.Second)

// 自定义重试条件
http.Client.AddRetryCondition(func(r *resty.Response, err error) bool {
    return r.StatusCode() == 429  // 仅在 429 状态码时重试
})
```

### 3. 日志记录

```go
// 开启调试模式（会打印请求和响应详情）
http.Client.SetDebug(true)

// 或在代码中记录
resp, err := http.Get(url)
logger.Info("http", "Request: %s, Status: %d", url, resp.StatusCode())
```

### 4. SSL 证书验证

```go
// 生产环境应该开启 SSL 验证
import "crypto/tls"

http.Client.SetTLSClientConfig(&tls.Config{
    InsecureSkipVerify: false,  // 开启验证
})
```

---

## 🔗 相关资源

- **Resty 官方文档**：https://github.com/go-resty/resty
- **HTTP 状态码**：https://httpstatuses.com/

---

**更新日期**：2025-10-16  
**版本**：v1.0
