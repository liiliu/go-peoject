package common

// ========================================
// 常量定义
// ========================================

// 缓存键前缀
const (
	// CacheTokenPrefix Token缓存前缀
	CacheTokenPrefix = "Token"
	// CacheUserPrefix 用户缓存前缀
	CacheUserPrefix = "User"
)

// Token相关
const (
	// TokenType Token类型
	TokenType = "Bearer"
)

// 日志标签
const (
	// LogTagApiError API错误
	LogTagApiError = "API错误"
	// LogTagRedisError Redis错误
	LogTagRedisError = "REDIS错误"
	// LogTagDBError 数据库错误
	LogTagDBError = "数据库错误"
	// LogTagAuthError 认证错误
	LogTagAuthError = "认证错误"
)

// 系统类型
const (
	// SysBackend 后台系统
	SysBackend = "backend"
	// SysWeb Web端
	SysWeb = "web"
	// SysApp App端
	SysApp = "app"
)
