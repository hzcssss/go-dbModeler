package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger 是应用的日志接口
type Logger struct {
	*zap.SugaredLogger
}

// New 创建一个新的日志实例
func New() *Logger {
	// 创建基本配置
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 创建日志文件
	logFile, _ := os.OpenFile("godbmodeler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// 同时输出到控制台和文件
	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config.EncoderConfig),
			zapcore.AddSync(os.Stdout),
			config.Level,
		),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			zapcore.AddSync(logFile),
			config.Level,
		),
	)

	// 创建logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()

	return &Logger{sugar}
}
