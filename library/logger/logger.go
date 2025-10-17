package logger

import (
	"os"
	"your_project/library/config"

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
		Compress:   true,                                    // 是否压缩
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
		// 生产模式：输出到控制台和文件，级别为Error
		AtomicLevel.SetLevel(zap.ErrorLevel)
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

// Info 信息日志（带标签）
func Info(tag, msg string, args ...interface{}) {
	if tag != "" {
		field := zap.Fields(zap.String("tag", tag))
		MyLogger.WithOptions(field).Sugar().Infof(msg, args...)
	} else {
		MyLogger.Sugar().Infof(msg, args...)
	}
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
