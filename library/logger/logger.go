package logger

import (
	"os"
	"your_project/library/config"

	"github.com/gofiber/fiber/v2"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MyLogger *zap.Logger
var AtomicLevel = zap.NewAtomicLevel()

// Init 初始化日志
func Init() {
	// 日志轮转配置
	hook := lumberjack.Logger{
		Filename:   config.GetConfig().Server.LogFileFolder, // 日志文件路径
		MaxSize:    128,                                     // 每个日志文件最大尺寸 (MB)
		MaxBackups: 30,                                      // 最多保留的备份数量
		MaxAge:     7,                                       // 文件最多保存天数
		Compress:   false,                                   // 是否压缩
		LocalTime:  true,                                    // 使用本地时间
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var core zapcore.Core
	if config.GetConfig().Server.Debug {
		// Debug模式：输出到控制台，级别为Debug
		AtomicLevel.SetLevel(zap.DebugLevel)
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			AtomicLevel,
		)
	} else {
		// 生产模式：输出到控制台和文件，级别为Info
		AtomicLevel.SetLevel(zap.InfoLevel)
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
			AtomicLevel,
		)
	}

	// 构造日志记录器
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	development := zap.Development()
	field := zap.Fields(zap.String("service_name", config.GetConfig().Server.Name))

	MyLogger = zap.New(core, caller, development, field, callerSkip)
}

// Debug 调试日志（带标签）
func Debug(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Debugf(msg, args...)
	} else {
		MyLogger.Sugar().Debugf(msg, args...)
	}
}

// Info 信息日志（带标签）
func Info(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Infof(msg, args...)
	} else {
		MyLogger.Sugar().Infof(msg, args...)
	}
}

// Warn 警告日志（带标签）
func Warn(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Warnf(msg, args...)
	} else {
		MyLogger.Sugar().Warnf(msg, args...)
	}
}

// Error 错误日志（带标签）
func Error(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Errorf(msg, args...)
	} else {
		MyLogger.Sugar().Errorf(msg, args...)
	}
}

// Fatal 致命错误日志（带标签）
func Fatal(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Fatalf(msg, args...)
	} else {
		MyLogger.Sugar().Fatalf(msg, args...)
	}
}

// GetTraceIDFromCtx 从Fiber Context中获取TraceID
func GetTraceIDFromCtx(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}
	if traceID, ok := c.Locals("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// DebugWithTrace 带TraceID的调试日志
func DebugWithTrace(c *fiber.Ctx, tag, msg string, args ...interface{}) {
	traceID := GetTraceIDFromCtx(c)
	fields := []zap.Field{}
	if tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if len(fields) > 0 {
		MyLogger.WithOptions(zap.Fields(fields...)).Sugar().Debugf(msg, args...)
	} else {
		MyLogger.Sugar().Debugf(msg, args...)
	}
}

// InfoWithTrace 带TraceID的信息日志
func InfoWithTrace(c *fiber.Ctx, tag, msg string, args ...interface{}) {
	traceID := GetTraceIDFromCtx(c)
	fields := []zap.Field{}
	if tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if len(fields) > 0 {
		MyLogger.WithOptions(zap.Fields(fields...)).Sugar().Infof(msg, args...)
	} else {
		MyLogger.Sugar().Infof(msg, args...)
	}
}

// WarnWithTrace 带TraceID的警告日志
func WarnWithTrace(c *fiber.Ctx, tag, msg string, args ...interface{}) {
	traceID := GetTraceIDFromCtx(c)
	fields := []zap.Field{}
	if tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if len(fields) > 0 {
		MyLogger.WithOptions(zap.Fields(fields...)).Sugar().Warnf(msg, args...)
	} else {
		MyLogger.Sugar().Warnf(msg, args...)
	}
}

// ErrorWithTrace 带TraceID的错误日志
func ErrorWithTrace(c *fiber.Ctx, tag, msg string, args ...interface{}) {
	traceID := GetTraceIDFromCtx(c)
	fields := []zap.Field{}
	if tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if len(fields) > 0 {
		MyLogger.WithOptions(zap.Fields(fields...)).Sugar().Errorf(msg, args...)
	} else {
		MyLogger.Sugar().Errorf(msg, args...)
	}
}

// FatalWithTrace 带TraceID的致命错误日志
func FatalWithTrace(c *fiber.Ctx, tag, msg string, args ...interface{}) {
	traceID := GetTraceIDFromCtx(c)
	fields := []zap.Field{}
	if tag != "" {
		fields = append(fields, zap.String("tag", tag))
	}
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if len(fields) > 0 {
		MyLogger.WithOptions(zap.Fields(fields...)).Sugar().Fatalf(msg, args...)
	} else {
		MyLogger.Sugar().Fatalf(msg, args...)
	}
}
