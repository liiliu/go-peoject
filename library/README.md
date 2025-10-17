# Library 工具库总览

## 📚 工具库列表

| 工具库 | 说明 | 文档 |
|--------|------|------|
| **http** | HTTP 客户端（基于 resty） | [使用指南](../docs/HTTP请求工具使用指南.md) |
| **excel** | Excel 文件读写、样式设置 | [使用指南](../docs/文件与数据处理工具使用指南.md) |
| **csv** | CSV 文件读写、数据转换 | [使用指南](../docs/文件与数据处理工具使用指南.md) |
| **file** | 文件下载、压缩、哈希操作 | [使用指南](../docs/文件与数据处理工具使用指南.md) |
| **util** | 常用工具函数集合 | [使用指南](../docs/工具函数使用指南.md) |
| **validator** | 参数验证 | [使用指南](../docs/参数验证使用指南.md) |
| **logger** | 日志记录 | - |
| **database** | 数据库操作（GORM） | - |
| **cache** | Redis 缓存 | - |
| **jwt** | JWT 认证 | - |

---

## 🚀 快速开始

### HTTP 请求

```go
import "your_project/library/http"

// GET 请求
resp, err := http.Get("https://api.example.com/users")

// POST JSON
data := map[string]interface{}{"name": "张三", "age": 25}
resp, err := http.Post("https://api.example.com/users", data)

// 带 Token
resp, err := http.Client.R().
    SetAuthToken("your-token").
    Get("https://api.example.com/protected")
```

### Excel 处理

```go
import "your_project/library/excel"

// 读取 Excel
data, err := excel.ReadSheetToMap("users.xlsx", "Sheet1")

// 导出 Excel（带样式）
headers := []string{"姓名", "年龄", "城市"}
data := [][]interface{}{
    {"张三", 25, "北京"},
    {"李四", 30, "上海"},
}
err := excel.ExportData("output.xlsx", "Sheet1", headers, data)
```

### CSV 处理

```go
import "your_project/library/csv"

// 读取 CSV
data, err := csv.ReadToMap("users.csv")

// 写入 CSV
headers := []string{"name", "age"}
users := []map[string]string{
    {"name": "张三", "age": "25"},
}
err := csv.WriteFromMap("output.csv", headers, users)
```

### 文件操作

```go
import "your_project/library/file"

// 下载文件
err := file.Download("https://example.com/file.pdf", "/path/to/save.pdf")

// 压缩文件
files := []string{"file1.txt", "file2.txt"}
err := file.ZipFiles("archive.zip", files)

// 计算 MD5
md5, err := file.MD5File("/path/to/file.txt")
```

### 工具函数

```go
import "your_project/library/util"

// 字符串操作
util.IsEmpty(" ")                        // true
util.MaskString("13812345678", 3, 4, "*") // "138****5678"

// 随机生成
code := util.GenerateCode(6)             // "892341"
orderNo := util.GenerateOrderNo()        // "1697456789123456"

// 验证
util.IsEmail("user@example.com")         // true
util.IsMobile("13812345678")             // true

// 类型转换
util.ToInt("123")                        // 123
util.ToString(456)                       // "456"
```

---

## 📦 依赖包

```toml
# go.mod
require (
    github.com/go-resty/resty/v2 v2.11.0      # HTTP 客户端
    github.com/xuri/excelize/v2 v2.8.1        # Excel 处理
    github.com/go-playground/validator/v10    # 参数验证
    # ... 其他依赖
)
```

安装依赖：
```bash
go mod tidy
```

---

## 📖 完整文档

- [HTTP 请求工具使用指南](../docs/HTTP请求工具使用指南.md)
- [文件与数据处理工具使用指南](../docs/文件与数据处理工具使用指南.md)
- [工具函数使用指南](../docs/工具函数使用指南.md)
- [参数验证使用指南](../docs/参数验证使用指南.md)

---

## 🎯 常见场景

### 场景 1：用户数据导入导出

```go
// 导入 Excel
data, _ := excel.ReadSheetToMap("users.xlsx", "Sheet1")
for _, record := range data {
    user := &model.User{
        Username: record["用户名"],
        Email:    record["邮箱"],
    }
    database.DB.Create(user)
}

// 导出 CSV
var users []model.User
database.DB.Find(&users)

headers := []string{"用户名", "邮箱", "注册时间"}
data := []map[string]string{}
for _, user := range users {
    data = append(data, map[string]string{
        "用户名":  user.Username,
        "邮箱":   user.Email,
        "注册时间": user.CreatedAt.Format("2006-01-02"),
    })
}
csv.WriteFromMap("users.csv", headers, data)
```

### 场景 2：调用第三方 API

```go
// 微信小程序登录
resp, err := http.Client.R().
    SetQueryParams(map[string]string{
        "appid":      "wx123456",
        "secret":     "secret",
        "js_code":    code,
        "grant_type": "authorization_code",
    }).
    Get("https://api.weixin.qq.com/sns/jscode2session")

var result WechatResponse
util.FromJSON(resp.String(), &result)
```

### 场景 3：文件上传下载

```go
// 上传处理
fileHeader, _ := c.FormFile("file")
savePath := "uploads/" + util.RandomString(16) + file.GetExtension(fileHeader.Filename)
c.SaveFile(fileHeader, savePath)

// 计算 MD5 防重
md5Hash, _ := file.MD5File(savePath)

// 下载远程文件
err := file.Download("https://example.com/file.pdf", "/local/path/file.pdf")
```

---

**项目**：Go Starter Template  
**更新日期**：2025-10-16
