package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

// InitLogger 初始化 zap logger
func InitLogger() {
	// 读取日志配置
	sec, err := Cfg.GetSection("log")
	if err != nil {
		// 如果配置不存在，使用默认配置
		Logger = getDefaultLogger()
		return
	}

	// 日志目录
	logDir := sec.Key("DIR").MustString("logs")
	// 日志级别: debug, info, warn, error
	levelStr := sec.Key("LEVEL").MustString("info")
	// 是否显示调用位置
	showLine := sec.Key("SHOW_LINE").MustBool(true)
	// 日志文件名
	logFileName := sec.Key("FILE_NAME").MustString("app.log")
	// 是否输出到控制台
	consoleOutput := sec.Key("CONSOLE_OUTPUT").MustBool(true)
	// 是否输出到文件
	fileOutput := sec.Key("FILE_OUTPUT").MustBool(true)
	// 日志文件最大大小(MB)
	maxSize := sec.Key("MAX_SIZE").MustInt(100)
	// 保留文件数量
	maxBackups := sec.Key("MAX_BACKUPS").MustInt(7)
	// 保留天数
	maxAge := sec.Key("MAX_AGE").MustInt(30)
	// 是否压缩
	compress := sec.Key("COMPRESS").MustBool(true)

	// 创建日志目录
	if fileOutput {
		// ModePerm:0777
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Printf("Failed to create log directory: %v\n", err)
			Logger = getDefaultLogger()
			return
		}
	}

	// 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		level = zapcore.InfoLevel
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
	}

	// 创建编码器
	// 编码器：决定日志格式化的方式，或者说输出格式
	var encoder zapcore.Encoder
	if RunMode == "debug" {
		// debug 模式使用 console 编码（彩色输出）
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 生产环境使用 JSON 编码
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建多个 core
	var cores []zapcore.Core

	// 控制台输出
	if consoleOutput {
		consoleCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 文件输出
	if fileOutput {
		logFile := filepath.Join(logDir, logFileName)
		fileWriter := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
			LocalTime:  true,
		}

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 文件统一使用 JSON
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 合并所有 core
	core := zapcore.NewTee(cores...)

	// 创建 logger
	Logger = zap.New(core)

	// 添加调用位置
	if showLine {
		Logger = Logger.WithOptions(zap.AddCaller())
	}

	// 添加堆栈跟踪（仅 error 及以上级别）
	Logger = Logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
}

// getDefaultLogger 获取默认 logger（当配置不存在时使用）
func getDefaultLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	return zap.New(core).WithOptions(zap.AddCaller())
}

// Sync 同步日志缓冲区（应用退出时调用）
func SyncLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
