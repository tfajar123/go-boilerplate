package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	_ = os.MkdirAll("logs", 0755)

	level := zapcore.InfoLevel
	if os.Getenv("APP_ENV") == "production" {
		level = zapcore.WarnLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    50,
			MaxBackups: 7,
			MaxAge:     30,
			Compress:   true,
		}),
		level, // ← dipakai di sini
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	Logger = zap.New(
		zapcore.NewTee(fileCore, consoleCore),
		zap.AddCaller(),
	)
}

func SyncLogger() {
	_ = Logger.Sync()
}
