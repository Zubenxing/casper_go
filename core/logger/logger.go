package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"casper_go/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// Init 初始化日志系统
func Init(cfg *config.LoggerConfig) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 创建日志目录
	logDir := filepath.Dir(cfg.GetLogPath())
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 配置日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.GetLogPath(),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	// 同时输出到文件和控制台
	mw := io.MultiWriter(os.Stdout, fileWriter)
	Log.SetOutput(mw)

	Log.Info("日志系统初始化完成")
	return nil
}

// Get 获取日志实例
func Get() *logrus.Logger {
	if Log == nil {
		Log = logrus.New()
	}
	return Log
}
