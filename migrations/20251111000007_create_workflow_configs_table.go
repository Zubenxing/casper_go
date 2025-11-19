package migrations

import (
	"casper_go/features/workflow"

	"gorm.io/gorm"
)

type CreateWorkflowConfigsTable20251111000007 struct{}

func (m *CreateWorkflowConfigsTable20251111000007) Version() string {
	return "20251111000007"
}

func (m *CreateWorkflowConfigsTable20251111000007) Name() string {
	return "create_workflow_configs_table"
}

func (m *CreateWorkflowConfigsTable20251111000007) Up(tx *gorm.DB) error {
	// 创建 workflow_configs 表（AutoMigrate 会自动创建索引）
	if err := tx.AutoMigrate(&workflow.WorkflowConfig{}); err != nil {
		return err
	}

	// 额外创建 enabled 字段的索引以提升查询性能
	if err := tx.Exec(`
		CREATE INDEX idx_workflow_configs_enabled 
		ON workflow_configs(enabled)
	`).Error; err != nil {
		// 索引可能已存在，忽略错误
	}

	return nil
}

func (m *CreateWorkflowConfigsTable20251111000007) Down(tx *gorm.DB) error {
	// 删除索引
	tx.Exec("DROP INDEX IF EXISTS idx_workflow_configs_enabled ON workflow_configs")

	// 删除表（会自动删除表上的所有索引）
	if err := tx.Migrator().DropTable(&workflow.WorkflowConfig{}); err != nil {
		return err
	}

	return nil
}
