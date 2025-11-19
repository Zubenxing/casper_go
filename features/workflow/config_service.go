package workflow

import (
	"casper_go/core/database"
	"errors"
)

// GetAllWorkflowConfigs 获取所有工作流配置
func GetAllWorkflowConfigs() ([]*WorkflowConfig, error) {
	var configs []*WorkflowConfig
	result := database.DB.Where("enabled = ?", true).Find(&configs)
	return configs, result.Error
}

// GetWorkflowConfigByWorkflowID 根据工作流ID获取配置
func GetWorkflowConfigByWorkflowID(workflowID string) (*WorkflowConfig, error) {
	var config WorkflowConfig
	result := database.DB.Where("workflow_id = ? AND enabled = ?", workflowID, true).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// GetWorkflowConfigByWorkflowName 根据工作流名称获取配置
func GetWorkflowConfigByWorkflowName(workflowName string) (*WorkflowConfig, error) {
	var config WorkflowConfig
	result := database.DB.Where("workflow_name = ? AND enabled = ?", workflowName, true).First(&config)
	if result.Error != nil {
		return nil, result.Error
	}
	return &config, nil
}

// CreateWorkflowConfig 创建工作流配置
func CreateWorkflowConfig(config *WorkflowConfig) error {
	// 检查是否已存在
	var existing WorkflowConfig
	result := database.DB.Where("workflow_id = ?", config.WorkflowID).First(&existing)
	if result.Error == nil {
		return errors.New("该工作流配置已存在")
	}

	return database.DB.Create(config).Error
}

// UpdateWorkflowConfig 更新工作流配置
func UpdateWorkflowConfig(id uint, updates *WorkflowConfig) error {
	return database.DB.Model(&WorkflowConfig{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWorkflowConfig 删除工作流配置
func DeleteWorkflowConfig(id uint) error {
	return database.DB.Delete(&WorkflowConfig{}, id).Error
}
