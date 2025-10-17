package http

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client 全局 HTTP 客户端
var Client *resty.Client

func init() {
	Client = NewClient()
}

// Config HTTP 客户端配置
type Config struct {
	Timeout       time.Duration     // 超时时间
	RetryCount    int               // 重试次数
	RetryWaitTime time.Duration     // 重试等待时间
	BaseURL       string            // 基础 URL
	Headers       map[string]string // 默认 Headers
	Debug         bool              // 是否开启调试模式
}

// NewClient 创建新的 HTTP 客户端
func NewClient() *resty.Client {
	client := resty.New()

	// 默认配置
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)

	// 设置默认 Header
	client.SetHeader("User-Agent", "Go-HTTP-Client/1.0")
	client.SetHeader("Accept", "application/json")

	// 跳过 SSL 验证（开发环境使用，生产环境建议关闭）
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return client
}

// NewClientWithConfig 使用自定义配置创建客户端
func NewClientWithConfig(cfg Config) *resty.Client {
	client := resty.New()

	if cfg.Timeout > 0 {
		client.SetTimeout(cfg.Timeout)
	} else {
		client.SetTimeout(30 * time.Second)
	}

	if cfg.RetryCount > 0 {
		client.SetRetryCount(cfg.RetryCount)
	}

	if cfg.RetryWaitTime > 0 {
		client.SetRetryWaitTime(cfg.RetryWaitTime)
	}

	if cfg.BaseURL != "" {
		client.SetBaseURL(cfg.BaseURL)
	}

	if cfg.Headers != nil {
		client.SetHeaders(cfg.Headers)
	}

	if cfg.Debug {
		client.SetDebug(true)
	}

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return client
}

// ========================================
// 便捷方法
// ========================================

// Get 发送 GET 请求
func Get(url string) (*resty.Response, error) {
	return Client.R().Get(url)
}

// GetWithParams 发送带参数的 GET 请求
func GetWithParams(url string, params map[string]string) (*resty.Response, error) {
	return Client.R().
		SetQueryParams(params).
		Get(url)
}

// GetWithHeaders 发送带自定义 Header 的 GET 请求
func GetWithHeaders(url string, headers map[string]string) (*resty.Response, error) {
	return Client.R().
		SetHeaders(headers).
		Get(url)
}

// Post 发送 POST 请求（JSON）
func Post(url string, body interface{}) (*resty.Response, error) {
	return Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
}

// PostForm 发送 POST 表单请求
func PostForm(url string, data map[string]string) (*resty.Response, error) {
	return Client.R().
		SetFormData(data).
		Post(url)
}

// PostJSON 发送 POST JSON 请求并解析响应
func PostJSON(url string, body interface{}, result interface{}) (*resty.Response, error) {
	return Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(result).
		Post(url)
}

func PostJSONWithHeaders(url string, body interface{}, headers map[string]string, result interface{}) (*resty.Response, error) {
	return Client.R().
		SetHeaders(headers).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(result).
		Post(url)
}

// Put 发送 PUT 请求
func Put(url string, body interface{}) (*resty.Response, error) {
	return Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Put(url)
}

// Delete 发送 DELETE 请求
func Delete(url string) (*resty.Response, error) {
	return Client.R().Delete(url)
}

// Patch 发送 PATCH 请求
func Patch(url string, body interface{}) (*resty.Response, error) {
	return Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Patch(url)
}

// ========================================
// 高级功能
// ========================================

// Request 通用请求构建器
type Request struct {
	client *resty.Client
	req    *resty.Request
}

// NewRequest 创建新的请求构建器
func NewRequest() *Request {
	return &Request{
		client: Client,
		req:    Client.R(),
	}
}

// SetURL 设置请求 URL
func (r *Request) SetURL(url string) *Request {
	r.req.URL = url
	return r
}

// SetMethod 设置请求方法
func (r *Request) SetMethod(method string) *Request {
	r.req.Method = method
	return r
}

// SetHeader 设置单个 Header
func (r *Request) SetHeader(key, value string) *Request {
	r.req.SetHeader(key, value)
	return r
}

// SetHeaders 设置多个 Headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	r.req.SetHeaders(headers)
	return r
}

// SetQueryParam 设置单个查询参数
func (r *Request) SetQueryParam(key, value string) *Request {
	r.req.SetQueryParam(key, value)
	return r
}

// SetQueryParams 设置多个查询参数
func (r *Request) SetQueryParams(params map[string]string) *Request {
	r.req.SetQueryParams(params)
	return r
}

// SetBody 设置请求体
func (r *Request) SetBody(body interface{}) *Request {
	r.req.SetBody(body)
	return r
}

// SetResult 设置响应结果对象
func (r *Request) SetResult(result interface{}) *Request {
	r.req.SetResult(result)
	return r
}

// SetError 设置错误响应对象
func (r *Request) SetError(err interface{}) *Request {
	r.req.SetError(err)
	return r
}

// SetBasicAuth 设置基础认证
func (r *Request) SetBasicAuth(username, password string) *Request {
	r.req.SetBasicAuth(username, password)
	return r
}

// SetAuthToken 设置 Bearer Token
func (r *Request) SetAuthToken(token string) *Request {
	r.req.SetAuthToken(token)
	return r
}

// Execute 执行请求
func (r *Request) Execute(method, url string) (*resty.Response, error) {
	return r.req.Execute(method, url)
}

// Send 发送请求（自动根据 Method 和 URL）
func (r *Request) Send() (*resty.Response, error) {
	return r.req.Send()
}

// ========================================
// 辅助函数
// ========================================

// IsSuccess 判断响应是否成功 (2xx)
func IsSuccess(resp *resty.Response) bool {
	return resp.StatusCode() >= 200 && resp.StatusCode() < 300
}

// GetString 获取响应字符串
func GetString(resp *resty.Response) string {
	return resp.String()
}

// GetBytes 获取响应字节数组
func GetBytes(resp *resty.Response) []byte {
	return resp.Body()
}

// GetStatusCode 获取状态码
func GetStatusCode(resp *resty.Response) int {
	return resp.StatusCode()
}

// ========================================
// 使用示例（注释）
// ========================================
/*
// 1. 简单的 GET 请求
resp, err := http.Get("https://api.example.com/users")
if err != nil {
    logger.Error("http", "Request failed: %v", err)
}
logger.Info("http", "Response: %s", resp.String())

// 2. 带参数的 GET 请求
params := map[string]string{
    "page": "1",
    "size": "10",
}
resp, err := http.GetWithParams("https://api.example.com/users", params)

// 3. POST JSON 请求
data := map[string]interface{}{
    "name": "张三",
    "age":  25,
}
resp, err := http.Post("https://api.example.com/users", data)

// 4. POST JSON 并解析响应
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
var result User
resp, err := http.PostJSON("https://api.example.com/users", data, &result)
logger.Info("http", "Created user ID: %d", result.ID)

// 5. 表单提交
formData := map[string]string{
    "username": "admin",
    "password": "123456",
}
resp, err := http.PostForm("https://api.example.com/login", formData)

// 6. 带 Token 的请求
req := http.NewRequest().
    SetAuthToken("your-token-here").
    SetQueryParam("status", "active").
    SetBody(data)
resp, err := req.Execute("POST", "https://api.example.com/orders")

// 7. 自定义超时（通过创建新客户端）
client := http.NewClient()
client.SetTimeout(5 * time.Second)
resp, err := client.R().Get("https://api.example.com/slow-endpoint")

// 或使用配置创建客户端
customClient := http.NewClientWithConfig(http.Config{
    Timeout:    10 * time.Second,
    RetryCount: 5,
    BaseURL:    "https://api.example.com",
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
    Debug: true,
})
resp, err = customClient.R().Get("/users")

// 8. 链式调用
resp, err := http.Client.R().
    SetHeader("Authorization", "Bearer token").
    SetQueryParam("status", "active").
    SetBody(map[string]string{"name": "test"}).
    Post("https://api.example.com/users")

// 9. 下载文件
resp, err := http.Client.R().
    SetOutput("/path/to/save/file.pdf").
    Get("https://example.com/file.pdf")

// 10. 上传文件
resp, err := http.Client.R().
    SetFile("file", "/path/to/file.jpg").
    SetFormData(map[string]string{
        "title": "My Photo",
    }).
    Post("https://api.example.com/upload")
*/
