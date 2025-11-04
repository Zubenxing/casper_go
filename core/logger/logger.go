package logger

import (
	"io"
	"os"
	"time"

	"casper_go/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局日志实例
var (
	Log         *logrus.Logger // 系统通用日志
	Certificate *logrus.Logger // 证书监控模块日志
	Password    *logrus.Logger // 密码管理模块日志
	Work        *logrus.Logger // 工作记录模块日志
)

// Init 初始化日志系统
func Init(cfg *config.LoggerConfig) error {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 解析日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}

	// 初始化系统通用日志
	Log = createLogger("logs/app.log", level, cfg)

	// 初始化证书模块日志
	Certificate = createLogger("logs/certificate.log", level, cfg)

	// 初始化密码模块日志
	Password = createLogger("logs/password.log", level, cfg)

	// 初始化工作记录模块日志
	Work = createLogger("logs/work.log", level, cfg)

	Log.Info("日志系统初始化完成")
	Log.Infof("已创建 4 个日志模块: app, certificate, password, work")

	return nil
}

// createLogger 创建一个日志实例
func createLogger(filename string, level logrus.Level, cfg *config.LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	// 设置日志级别
	logger.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 配置日志轮转
	fileWriter := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,    // 每个日志文件最大尺寸（MB）
		MaxBackups: cfg.MaxBackups, // 保留的旧日志文件数量
		MaxAge:     cfg.MaxAge,     // 保留天数
		Compress:   cfg.Compress,   // 是否压缩
		LocalTime:  true,
	}

	// 同时输出到文件和控制台
	mw := io.MultiWriter(os.Stdout, fileWriter)
	logger.SetOutput(mw)

	return logger
}

// Get 获取系统通用日志实例
func Get() *logrus.Logger {
	if Log == nil {
		Log = logrus.New()
	}
	return Log
}

// GetCertificate 获取证书模块日志实例
func GetCertificate() *logrus.Logger {
	if Certificate == nil {
		Certificate = logrus.New()
	}
	return Certificate
}

// GetPassword 获取密码模块日志实例
func GetPassword() *logrus.Logger {
	if Password == nil {
		Password = logrus.New()
	}
	return Password
}
