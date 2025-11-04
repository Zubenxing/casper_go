package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"casper_go/config"
	"casper_go/core/database"
	"casper_go/core/logger"
	_ "casper_go/docs"
	"casper_go/features/auth"
	"casper_go/features/certificate"
	"casper_go/features/password"
	"casper_go/migrations"
	"casper_go/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 检查是否是迁移命令
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}

	// 加载配置
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Logger); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("=== Casper Platform 启动中 ===")
	logger.Log.Infof("版本: %s", cfg.Site.Version)
	logger.Log.Infof("模式: %s", cfg.Site.Mode)

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Log.Fatalf("初始化数据库失败: %v", err)
	}

	// 数据库迁移（根据配置）
	if cfg.DatabaseInit.AutoMigrate {
		if err := database.DB.AutoMigrate(
			&auth.User{},
			&auth.UserToken{},
			&certificate.Monitor{},
			&password.Account{},
		); err != nil {
			logger.Log.Fatalf("数据库迁移失败: %v", err)
		}
		logger.Log.Info("数据库迁移完成")
	} else {
		logger.Log.Info("数据库自动迁移已禁用")
	}

	// 初始化默认数据（根据配置）
	if cfg.DatabaseInit.InitDefaultData {
		if err := auth.InitDefaultData(); err != nil {
			logger.Log.Fatalf("初始化默认数据失败: %v", err)
		}
		logger.Log.Info("默认数据初始化完成")
	} else {
		logger.Log.Info("默认数据初始化已禁用")
	}

	// 设置 Gin 模式
	if cfg.Site.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置路由
	r := router.Setup(cfg.Site.Mode, cfg)

	// 启动定时任务
	go startScheduledTasks(cfg)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Site.Host, cfg.Site.Port)
	logger.Log.Infof("服务器启动在: %s", addr)

	apiPath := "/api"
	if cfg.API.EnableVersion {
		apiPath = "/api/" + cfg.API.CurrentVersion
	}
	logger.Log.Infof("API 路径: %s", apiPath)
	logger.Log.Infof("Swagger 文档: http://localhost:%d/swagger/index.html", cfg.Site.Port)

	if cfg.API.RateLimitEnabled {
		logger.Log.Infof("API 限流已启用: %d 请求/分钟", cfg.API.RateLimitRequests)
	}
	if cfg.API.RequestIDEnabled {
		logger.Log.Info("请求ID追踪已启用")
	}

	// 优雅关闭
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("正在关闭服务器...")

	// 关闭数据库连接
	sqlDB, err := database.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	logger.Log.Info("服务器已安全关闭")
}

// startScheduledTasks 启动定时任务
func startScheduledTasks(cfg *config.Config) {
	ticker := time.NewTicker(time.Duration(cfg.Certificate.CheckInterval) * time.Second)
	defer ticker.Stop()

	logger.Log.Info("证书检查定时任务已启动")

	for {
		select {
		case <-ticker.C:
			logger.Log.Info("开始执行定时证书检查...")
			result, err := certificate.CheckAll()
			if err != nil {
				logger.Log.Errorf("证书检查失败: %v", err)
			} else {
				logger.Log.Infof("证书检查完成: 总数=%d, 成功=%d, 失败=%d, 耗时=%s",
					result.Total, result.Success, result.Failed, result.Duration)
			}
		}
	}
}

// runMigrations 运行数据库迁移
func runMigrations() {
	// 加载配置
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Logger); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("=== 开始数据库迁移 ===")

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建迁移管理器
	migrator := database.NewMigrator(database.DB)

	// 注册所有迁移
	allMigrations := migrations.GetAllMigrations()
	for _, m := range allMigrations {
		migrator.Register(m)
	}

	// 显示迁移状态
	if err := migrator.Status(); err != nil {
		logger.Log.Fatalf("获取迁移状态失败: %v", err)
	}

	logger.Log.Info("========================================")
	logger.Log.Info("开始执行迁移...")

	// 执行迁移
	if err := migrator.Up(); err != nil {
		logger.Log.Fatalf("迁移失败: %v", err)
	}

	logger.Log.Info("✅ 所有迁移已完成")
	os.Exit(0)
}
