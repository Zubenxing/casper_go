package migrations

import (
	"casper_go/core/database"
	"casper_go/core/logger"
)

// AddIssueImagesContent20251104000006 为问题记录添加图片和富文本内容字段
type AddIssueImagesContent20251104000006 struct{}

func (m *AddIssueImagesContent20251104000006) Version() string {
	return "20251104000006"
}

func (m *AddIssueImagesContent20251104000006) Name() string {
	return "为问题记录添加图片和富文本内容字段"
}

func (m *AddIssueImagesContent20251104000006) Up() error {
	logger.Log.Info("开始执行迁移: 为问题记录添加图片和富文本内容字段")

	// 添加 images 字段（JSON 数组，存储图片 URL）
	if err := database.DB.Exec(`
		ALTER TABLE work_issues 
		ADD COLUMN IF NOT EXISTS images JSON COMMENT '问题截图(最多3张)' AFTER description
	`).Error; err != nil {
		logger.Log.Errorf("添加 images 字段失败: %v", err)
		return err
	}

	// 添加 content 字段（富文本内容）
	if err := database.DB.Exec(`
		ALTER TABLE work_issues 
		ADD COLUMN IF NOT EXISTS content TEXT COMMENT '问题详细内容(富文本)' AFTER images
	`).Error; err != nil {
		logger.Log.Errorf("添加 content 字段失败: %v", err)
		return err
	}

	logger.Log.Info("迁移完成: 问题记录图片和富文本字段添加成功")
	return nil
}

func (m *AddIssueImagesContent20251104000006) Down() error {
	return nil // 不支持回滚
}
