# HTTP 客户端工具

基于 `go-resty/resty` 的 HTTP 客户端封装。

## 快速开始

```go
import "your_project/library/http"

// 简单 GET 请求
resp, err := http.Get("https://api.example.com/users")

// POST JSON
data := map[string]interface{}{
    "name": "张三",
    "age": 25,
}
resp, err := http.Post("https://api.example.com/users", data)

// 带 Token 的请求
resp, err := http.Client.R().
    SetAuthToken("your-jwt-token").
    Get("https://api.example.com/protected")
```

## 常用方法

| 方法 | 说明 |
|------|------|
| `Get(url)` | 发送 GET 请求 |
| `Post(url, body)` | 发送 POST JSON 请求 |
| `PostForm(url, data)` | 发送 POST 表单请求 |
| `PostJSON(url, body, result)` | POST 并解析响应 |
| `Put(url, body)` | 发送 PUT 请求 |
| `Delete(url)` | 发送 DELETE 请求 |

## 超时设置

**重要**：resty 的超时需要在 Client 级别设置，不能在单个请求上设置。

```go
// 方式 1：创建新客户端
client := http.NewClient()
client.SetTimeout(5 * time.Second)
resp, err := client.R().Get(url)

// 方式 2：使用配置
client := http.NewClientWithConfig(http.Config{
    Timeout: 5 * time.Second,
})

// 方式 3：修改全局客户端（影响所有后续请求）
http.Client.SetTimeout(5 * time.Second)
```

## 链式调用

```go
resp, err := http.Client.R().
    SetAuthToken("token").
    SetQueryParam("page", "1").
    SetHeader("Custom-Header", "value").
    SetBody(data).
    Post("https://api.example.com/users")
```

## 完整文档

查看 `docs/HTTP请求工具使用指南.md` 获取详细文档和更多示例。
