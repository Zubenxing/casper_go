package migrations

import (
	"casper_go/features/work"

	"gorm.io/gorm"
)

type CreateWorksTables20251103000004 struct{}

func (m *CreateWorksTables20251103000004) Version() string {
	return "20251103000004"
}

func (m *CreateWorksTables20251103000004) Name() string {
	return "create_works_tables"
}

func (m *CreateWorksTables20251103000004) Up(tx *gorm.DB) error {
	// 创建 works 表
	if err := tx.AutoMigrate(&work.Work{}); err != nil {
		return err
	}

	// 创建 work_issues 表
	if err := tx.AutoMigrate(&work.WorkIssue{}); err != nil {
		return err
	}

	// 创建额外的复合索引
	if err := tx.Exec(`
		CREATE INDEX idx_work_user_status_created 
		ON works(user_id, status, created_at DESC)
	`).Error; err != nil {
		// 索引可能已存在，忽略错误
		return nil
	}

	if err := tx.Exec(`
		CREATE INDEX idx_issue_user_status_created 
		ON work_issues(user_id, status, created_at DESC)
	`).Error; err != nil {
		// 索引可能已存在，忽略错误
		return nil
	}

	return nil
}

func (m *CreateWorksTables20251103000004) Down(tx *gorm.DB) error {
	// 删除索引
	tx.Exec("DROP INDEX IF EXISTS idx_work_user_status_created ON works")
	tx.Exec("DROP INDEX IF EXISTS idx_issue_user_status_created ON work_issues")

	// 删除表
	if err := tx.Migrator().DropTable(&work.WorkIssue{}); err != nil {
		return err
	}

	if err := tx.Migrator().DropTable(&work.Work{}); err != nil {
		return err
	}

	return nil
}
