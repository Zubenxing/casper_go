package main

import (
	"flag"
	"fmt"
	"time"

	"casper_go/config"
	"casper_go/core/database"
	"casper_go/core/logger"
)

func main() {
	configPath := flag.String("config", "config", "配置文件目录路径")
	flag.Parse()

	// 初始化日志
	logger.Init(&config.LoggerConfig{
		Level:      "info",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Log.Fatalf("配置加载失败: %v", err)
	}
	config.GlobalConfig = cfg

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Log.Fatalf("数据库初始化失败: %v", err)
	}

	fmt.Println("=== 数据库性能检查工具 ===\n")

	// 1. 检查 user_tokens 索引
	fmt.Println("📊 1. 检查 user_tokens 表索引:")
	checkIndexes("user_tokens")

	// 2. 检查 password_accounts 索引
	fmt.Println("\n📊 2. 检查 password_accounts 表索引:")
	checkIndexes("password_accounts")

	// 3. 测试 JWT 验证查询性能
	fmt.Println("\n⏱️  3. JWT 验证查询性能测试:")
	testQuery(`SELECT count(*) FROM user_tokens 
		WHERE token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiaXNzIjoiY2FzcGVyX3BsYXRmb3JtIiwiZXhwIjoxNzYyMjIxMTEwLCJpYXQiOjE3NjIxMzQ3MTB9.JWLzgHv2_iEucDs-xSTvpRSSqzAk208eHJqRxZ-Mrxg' 
		AND is_valid = true 
		AND expires_at > NOW()`)

	// 4. 测试密码账户查询性能
	fmt.Println("\n⏱️  4. 密码账户列表查询性能测试:")
	testQuery(`SELECT count(*) FROM password_accounts 
		WHERE user_id = 1 
		AND deleted_at IS NULL`)

	// 5. 测试带排序的查询
	fmt.Println("\n⏱️  5. 带排序的密码账户查询性能测试:")
	testQuery(`SELECT * FROM password_accounts 
		WHERE user_id = 1 
		AND deleted_at IS NULL 
		ORDER BY is_fav DESC, created_at ASC 
		LIMIT 20`)

	// 6. 分析查询计划
	fmt.Println("\n🔍 6. JWT 验证查询计划分析:")
	explainQuery(`SELECT count(*) FROM user_tokens 
		WHERE token = 'test_token' 
		AND is_valid = true 
		AND expires_at > NOW()`)

	fmt.Println("\n🔍 7. 密码账户查询计划分析:")
	explainQuery(`SELECT * FROM password_accounts 
		WHERE user_id = 1 
		AND deleted_at IS NULL 
		ORDER BY is_fav DESC, created_at ASC 
		LIMIT 20`)

	fmt.Println("\n✅ 检查完成！")
}

func checkIndexes(tableName string) {
	type Index struct {
		Table       string `gorm:"column:Table"`
		NonUnique   int    `gorm:"column:Non_unique"`
		KeyName     string `gorm:"column:Key_name"`
		SeqInIndex  int    `gorm:"column:Seq_in_index"`
		ColumnName  string `gorm:"column:Column_name"`
		Collation   string `gorm:"column:Collation"`
		Cardinality int64  `gorm:"column:Cardinality"`
		IndexType   string `gorm:"column:Index_type"`
	}

	var indexes []Index
	if err := database.DB.Raw("SHOW INDEX FROM " + tableName).Scan(&indexes).Error; err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		return
	}

	if len(indexes) == 0 {
		fmt.Println("⚠️  没有找到任何索引！")
		return
	}

	currentKey := ""
	for _, idx := range indexes {
		if idx.KeyName != currentKey {
			fmt.Printf("  ✓ %s (%s)\n", idx.KeyName, idx.IndexType)
			currentKey = idx.KeyName
		}
		fmt.Printf("    └─ %s (cardinality: %d)\n", idx.ColumnName, idx.Cardinality)
	}
}

func testQuery(query string) {
	start := time.Now()
	var count int64
	err := database.DB.Raw(query).Scan(&count).Error
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		return
	}

	if elapsed > 100*time.Millisecond {
		fmt.Printf("⚠️  耗时: %v (慢查询)\n", elapsed)
	} else {
		fmt.Printf("✅ 耗时: %v\n", elapsed)
	}
}

func explainQuery(query string) {
	type ExplainResult struct {
		ID           int     `gorm:"column:id"`
		SelectType   string  `gorm:"column:select_type"`
		Table        string  `gorm:"column:table"`
		Type         string  `gorm:"column:type"`
		PossibleKeys *string `gorm:"column:possible_keys"`
		Key          *string `gorm:"column:key"`
		KeyLen       *string `gorm:"column:key_len"`
		Ref          *string `gorm:"column:ref"`
		Rows         int64   `gorm:"column:rows"`
		Extra        *string `gorm:"column:Extra"`
	}

	var results []ExplainResult
	if err := database.DB.Raw("EXPLAIN " + query).Scan(&results).Error; err != nil {
		fmt.Printf("❌ 分析失败: %v\n", err)
		return
	}

	for _, r := range results {
		fmt.Printf("  表: %s\n", r.Table)
		fmt.Printf("  类型: %s\n", r.Type)
		if r.PossibleKeys != nil {
			fmt.Printf("  可用索引: %s\n", *r.PossibleKeys)
		} else {
			fmt.Println("  可用索引: NULL")
		}
		if r.Key != nil {
			fmt.Printf("  使用索引: %s ✅\n", *r.Key)
		} else {
			fmt.Println("  使用索引: NULL ❌ (全表扫描)")
		}
		fmt.Printf("  扫描行数: %d\n", r.Rows)
		if r.Extra != nil {
			fmt.Printf("  额外信息: %s\n", *r.Extra)
		}
	}
}
