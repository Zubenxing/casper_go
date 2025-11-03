package main

import (
	"flag"
	"fmt"
	"os"

	"casper_go/config"
	"casper_go/core/database"
	"casper_go/core/logger"
	"casper_go/migrations"
)

func main() {
	// 定义命令行参数
	action := flag.String("action", "up", "迁移操作: up(执行迁移), down(回滚), status(查看状态)")
	configPath := flag.String("config", "config", "配置文件目录路径")
	flag.Parse()

	// 初始化日志
	if err := logger.Init(&config.LoggerConfig{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	}); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	logger.Log.Info("=== Casper 数据库迁移工具 ===")

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Log.Fatalf("加载配置失败: %v", err)
	}

	// 设置全局配置
	config.GlobalConfig = cfg

	// 初始化数据库连接
	if err = database.Init(&config.GlobalConfig.Database); err != nil {
		logger.Log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建迁移管理器
	migrator := database.NewMigrator(database.DB)

	// 注册所有迁移
	allMigrations := migrations.GetAllMigrations()
	for _, m := range allMigrations {
		migrator.Register(m)
	}

	logger.Log.Infof("已注册 %d 个迁移文件", len(allMigrations))

	// 执行操作
	switch *action {
	case "up":
		logger.Log.Info("开始执行迁移...")
		err = migrator.Up()
	case "down":
		logger.Log.Info("开始回滚迁移...")
		err = migrator.Down()
	case "status":
		logger.Log.Info("查看迁移状态...")
		err = migrator.Status()
	default:
		logger.Log.Fatalf("未知操作: %s (支持: up, down, status)", *action)
	}

	if err != nil {
		logger.Log.Fatalf("操作失败: %v", err)
	}

	logger.Log.Info("=== 完成 ===")
}
