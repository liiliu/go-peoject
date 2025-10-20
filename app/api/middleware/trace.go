package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TraceIDKey 是context中存储TraceID的key
const TraceIDKey = "trace_id"

// TraceID 中间件，用于生成或提取TraceID并存储到context中
func TraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 尝试从请求头中获取TraceID
		traceID := c.Get("X-Trace-ID")
		
		// 如果请求头中没有，尝试使用Fiber自带的RequestID
		if traceID == "" {
			traceID = c.GetRespHeader(fiber.HeaderXRequestID)
		}
		
		// 如果还是没有，生成一个新的UUID
		if traceID == "" {
			traceID = uuid.New().String()
		}
		
		// 将TraceID存储到context的Locals中
		c.Locals(TraceIDKey, traceID)
		
		// 将TraceID添加到响应头中
		c.Set("X-Trace-ID", traceID)
		
		return c.Next()
	}
}

// GetTraceID 从context中获取TraceID
func GetTraceID(c *fiber.Ctx) string {
	if traceID, ok := c.Locals(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
