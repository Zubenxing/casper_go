package database

import (
	"fmt"
	"sort"
	"time"

	"casper_go/core/logger"

	"gorm.io/gorm"
)

// Migration 迁移接口
type Migration interface {
	// Version 返回迁移版本号（格式：YYYYMMDDHHMMSS）
	Version() string
	// Name 返回迁移名称
	Name() string
	// Up 执行迁移
	Up(db *gorm.DB) error
	// Down 回滚迁移
	Down(db *gorm.DB) error
}

// MigrationHistory 迁移历史记录
type MigrationHistory struct {
	ID        uint      `gorm:"primarykey;column:id;comment:Primary ID"`
	Version   string    `gorm:"type:varchar(20);column:version;uniqueIndex;not null;comment:迁移版本号"`
	Migration string    `gorm:"type:varchar(255);column:migration;not null;comment:迁移名称"`
	Batch     int       `gorm:"column:batch;default:0;not null;comment:批次号"`
	CreatedAt time.Time `gorm:"column:created_at;comment:执行时间"`
}

// TableName 指定表名
func (MigrationHistory) TableName() string {
	return "migrations"
}

// Migrator 迁移管理器
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator 创建迁移管理器
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register 注册迁移
func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// ensureHistoryTable 确保迁移历史表存在
func (m *Migrator) ensureHistoryTable() error {
	return m.db.AutoMigrate(&MigrationHistory{})
}

// getExecutedMigrations 获取已执行的迁移版本
func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	var histories []MigrationHistory
	if err := m.db.Find(&histories).Error; err != nil {
		return nil, err
	}

	executed := make(map[string]bool)
	for _, h := range histories {
		executed[h.Version] = true
	}
	return executed, nil
}

// sortMigrations 按版本号排序迁移
func (m *Migrator) sortMigrations() {
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version() < m.migrations[j].Version()
	})
}

// getNextBatch 获取下一个批次号
func (m *Migrator) getNextBatch() (int, error) {
	var maxBatch int
	if err := m.db.Model(&MigrationHistory{}).Select("COALESCE(MAX(batch), 0)").Scan(&maxBatch).Error; err != nil {
		return 0, err
	}
	return maxBatch + 1, nil
}

// Up 执行所有待执行的迁移
func (m *Migrator) Up() error {
	if err := m.ensureHistoryTable(); err != nil {
		return fmt.Errorf("创建迁移历史表失败: %w", err)
	}

	executed, err := m.getExecutedMigrations()
	if err != nil {
		return fmt.Errorf("获取迁移历史失败: %w", err)
	}

	m.sortMigrations()

	// 收集需要执行的迁移
	var pendingMigrations []Migration
	for _, migration := range m.migrations {
		if !executed[migration.Version()] {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	if len(pendingMigrations) == 0 {
		logger.Log.Info("数据库已是最新，无需迁移")
		return nil
	}

	// 获取新批次号
	batch, err := m.getNextBatch()
	if err != nil {
		return fmt.Errorf("获取批次号失败: %w", err)
	}

	// 执行待迁移
	for _, migration := range pendingMigrations {
		version := migration.Version()
		name := migration.Name()

		logger.Log.Infof("执行迁移: %s - %s (batch: %d)", version, name, batch)

		// 在事务中执行迁移
		if err := m.db.Transaction(func(tx *gorm.DB) error {
			// 执行迁移
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("迁移失败: %w", err)
			}

			// 记录迁移历史
			history := &MigrationHistory{
				Version:   version,
				Migration: name,
				Batch:     batch,
				CreatedAt: time.Now(),
			}
			if err := tx.Create(history).Error; err != nil {
				return fmt.Errorf("记录迁移历史失败: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("迁移 %s 失败: %w", version, err)
		}

		logger.Log.Infof("✓ 迁移完成: %s", version)
	}

	logger.Log.Infof("所有迁移已完成 (共 %d 个, batch: %d)", len(pendingMigrations), batch)
	return nil
}

// Down 回滚最后一个批次的迁移
func (m *Migrator) Down() error {
	if err := m.ensureHistoryTable(); err != nil {
		return fmt.Errorf("创建迁移历史表失败: %w", err)
	}

	// 获取最后一个批次号
	var maxBatch int
	if err := m.db.Model(&MigrationHistory{}).Select("MAX(batch)").Scan(&maxBatch).Error; err != nil {
		return fmt.Errorf("获取批次号失败: %w", err)
	}

	if maxBatch == 0 {
		logger.Log.Info("没有可回滚的迁移")
		return nil
	}

	// 获取该批次的所有迁移（按执行顺序倒序回滚）
	var histories []MigrationHistory
	if err := m.db.Where("batch = ?", maxBatch).Order("id DESC").Find(&histories).Error; err != nil {
		return fmt.Errorf("获取迁移历史失败: %w", err)
	}

	if len(histories) == 0 {
		logger.Log.Info("没有可回滚的迁移")
		return nil
	}

	logger.Log.Infof("开始回滚 batch %d (共 %d 个迁移)", maxBatch, len(histories))

	// 回滚每个迁移
	for _, history := range histories {
		// 查找对应的迁移
		var targetMigration Migration
		for _, migration := range m.migrations {
			if migration.Version() == history.Version {
				targetMigration = migration
				break
			}
		}

		if targetMigration == nil {
			logger.Log.Warnf("找不到迁移定义: %s，跳过", history.Version)
			continue
		}

		logger.Log.Infof("回滚迁移: %s - %s", history.Version, history.Migration)

		// 在事务中回滚迁移
		if err := m.db.Transaction(func(tx *gorm.DB) error {
			// 执行回滚
			if err := targetMigration.Down(tx); err != nil {
				return fmt.Errorf("回滚失败: %w", err)
			}

			// 删除迁移历史记录
			if err := tx.Delete(&history).Error; err != nil {
				return fmt.Errorf("删除迁移历史失败: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("回滚 %s 失败: %w", history.Version, err)
		}

		logger.Log.Infof("✓ 回滚完成: %s", history.Version)
	}

	logger.Log.Infof("批次 %d 回滚完成", maxBatch)
	return nil
}

// Status 显示迁移状态
func (m *Migrator) Status() error {
	if err := m.ensureHistoryTable(); err != nil {
		return fmt.Errorf("创建迁移历史表失败: %w", err)
	}

	// 获取所有已执行的迁移记录（带批次信息）
	var histories []MigrationHistory
	if err := m.db.Order("batch ASC, id ASC").Find(&histories).Error; err != nil {
		return fmt.Errorf("获取迁移历史失败: %w", err)
	}

	// 构建已执行迁移的 map
	executed := make(map[string]*MigrationHistory)
	for i := range histories {
		executed[histories[i].Version] = &histories[i]
	}

	m.sortMigrations()

	logger.Log.Info("迁移状态:")
	logger.Log.Info("========================================")

	if len(histories) > 0 {
		logger.Log.Info("已执行的迁移:")
		currentBatch := 0
		for _, history := range histories {
			if history.Batch != currentBatch {
				currentBatch = history.Batch
				logger.Log.Infof("  Batch %d:", currentBatch)
			}
			logger.Log.Infof("    [✓] %s - %s", history.Version, history.Migration)
		}
		logger.Log.Info("----------------------------------------")
	}

	// 显示待执行的迁移
	var pendingMigrations []Migration
	for _, migration := range m.migrations {
		if _, exists := executed[migration.Version()]; !exists {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	if len(pendingMigrations) > 0 {
		logger.Log.Info("待执行的迁移:")
		for _, migration := range pendingMigrations {
			logger.Log.Infof("    [ ] %s - %s", migration.Version(), migration.Name())
		}
	} else {
		logger.Log.Info("✓ 数据库已是最新，无待执行迁移")
	}

	logger.Log.Info("========================================")
	logger.Log.Infof("统计: 已执行 %d 个, 待执行 %d 个", len(histories), len(pendingMigrations))
	return nil
}
