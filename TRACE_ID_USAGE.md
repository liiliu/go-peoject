# TraceID 使用指南

## 概述
项目已支持完整的 TraceID 链路追踪功能，包括：
- ✅ 日志中记录 traceId
- ✅ 接口响应中返回 trace_id
- ✅ 自动生成和传递 traceId

## 功能说明

### 1. TraceID 自动生成
系统会自动为每个请求生成唯一的 TraceID：
- 优先从请求头 `X-Trace-ID` 中获取
- 如果没有，使用 Fiber 的 RequestID
- 都没有时，自动生成 UUID

### 2. 响应结构
所有接口响应都会包含 `trace_id` 字段：

```json
{
  "result": 1,
  "code": 1000,
  "msg": "成功",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "data": {
    "user": "张三"
  }
}
```

### 3. 日志记录
所有日志都会自动记录 traceId，方便链路追踪。

## 使用方法

### 方式一：使用便捷方法（最推荐）⭐

自动从 Context 获取 traceId，代码最简洁：

```go
package example

import (
    "your_project/app/view"
    "your_project/library/logger"
    "github.com/gofiber/fiber/v2"
)

func ExampleHandler(c *fiber.Ctx) error {
    // 1. 使用带 TraceID 的日志
    logger.InfoWithTrace(c, "example", "处理请求开始")
    
    // 2. 业务逻辑
    data := fiber.Map{
        "user": "张三",
        "age":  25,
    }
    
    // 3. 使用便捷方法返回，自动带上 traceId
    return c.JSON(view.SuccessWithCtx(c, data))
}

// 错误响应示例
func ExampleErrorHandler(c *fiber.Ctx) error {
    logger.ErrorWithTrace(c, "example", "发生错误: %v", "参数验证失败")
    return c.JSON(view.ErrorWithCtx(c, view.CodeInvalidParams))
}

// 自定义错误消息示例
func ExampleCustomErrorHandler(c *fiber.Ctx) error {
    logger.WarnWithTrace(c, "example", "自定义错误提示")
    return c.JSON(view.ErrorWithMsgCtx(c, view.CodeFailed, "自定义的错误信息"))
}
```

### 方式二：手动传递 TraceID

如果需要手动控制 traceId：

```go
import "your_project/app/api/middleware"

func ExampleHandler(c *fiber.Ctx) error {
    // 手动获取 traceId
    traceId := middleware.GetTraceID(c)
    
    logger.InfoWithTrace(c, "example", "处理请求")
    
    // 手动传递 traceId
    return c.JSON(view.SuccessResult(data, traceId))
}
```

## 日志函数说明

### 带 TraceID 的日志函数
- `logger.InfoWithTrace(c, tag, msg, args...)`   - 信息日志
- `logger.DebugWithTrace(c, tag, msg, args...)`  - 调试日志（仅 debug=true 时记录）
- `logger.WarnWithTrace(c, tag, msg, args...)`   - 警告日志
- `logger.ErrorWithTrace(c, tag, msg, args...)`  - 错误日志
- `logger.FatalWithTrace(c, tag, msg, args...)`  - 致命错误日志

### 传统日志函数（仍可使用）
- `logger.Info(tag, msg, args...)`
- `logger.Debug(tag, msg, args...)`
- `logger.Warn(tag, msg, args...)`
- `logger.Error(tag, msg, args...)`
- `logger.Fatal(tag, msg, args...)`

### 日志级别说明
- **Debug模式** (`debug = true`)：记录 Debug 及以上所有日志，输出到控制台
- **生产模式** (`debug = false`)：记录 Info 及以上日志，输出到控制台和文件

## 响应函数说明

### 便捷方法（推荐）⭐

自动从 Context 获取 traceId：

```go
// 成功响应（有数据）
view.SuccessWithCtx(c, data)

// 成功响应（无数据）
view.SuccessWithoutDataCtx(c)

// 错误响应
view.ErrorWithCtx(c, code)

// 自定义错误消息
view.ErrorWithMsgCtx(c, code, msg)
```

### 传统方法

支持可选的 traceId 参数：

```go
// 成功响应（有数据）
view.SuccessResult(data, traceId)
view.SuccessResult(data) // 不带 traceId 也可以

// 成功响应（无数据）
view.SuccessResultWithoutData(traceId)
view.SuccessResultWithoutData()

// 错误响应
view.ErrorResult(code, traceId)
view.ErrorResult(code)

// 自定义错误消息
view.ErrorResultWithMsg(code, msg, traceId)
view.ErrorResultWithMsg(code, msg)
```

## 客户端使用

### 主动传递 TraceID
客户端可以在请求头中传递 TraceID：

```bash
curl -H "X-Trace-ID: my-custom-trace-id" http://localhost:8080/v1/health
```

### 获取响应中的 TraceID
从响应体的 `trace_id` 字段或响应头的 `X-Trace-ID` 中获取。

## 日志输出示例

```json
{
  "level": "info",
  "time": "2025-10-20T13:30:00.000+0800",
  "caller": "example/handler.go:15",
  "msg": "处理请求开始",
  "service_name": "your_project",
  "tag": "example",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 注意事项

1. **TraceID 中间件已自动注册**：无需手动添加
2. **建议使用带 Trace 的日志函数**：便于日志关联查询
3. **响应中推荐携带 traceId**：方便前端排查问题
4. **traceId 参数是可选的**：保持向后兼容性
