package workflow

import (
	"casper_go/core/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetWorkflowConfigsAPI 获取所有工作流配置
func GetWorkflowConfigsAPI(c *gin.Context) {
	configs, err := GetAllWorkflowConfigs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取配置失败: "+err.Error())
		return
	}

	// 转换为 DTO
	dtos := make([]*WorkflowConfigDTO, len(configs))
	for i, config := range configs {
		dtos[i] = config.ToDTO()
	}

	response.Success(c, dtos)
}

// GetWorkflowConfigByWorkflowIDAPI 根据工作流ID获取配置
func GetWorkflowConfigByWorkflowIDAPI(c *gin.Context) {
	workflowID := c.Param("id")
	if workflowID == "" {
		response.Error(c, http.StatusBadRequest, "工作流ID不能为空")
		return
	}

	config, err := GetWorkflowConfigByWorkflowID(workflowID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "配置不存在")
		return
	}

	response.Success(c, config.ToDTO())
}

// CreateWorkflowConfigAPI 创建工作流配置
func CreateWorkflowConfigAPI(c *gin.Context) {
	var config WorkflowConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := CreateWorkflowConfig(&config); err != nil {
		response.Error(c, http.StatusInternalServerError, "创建配置失败: "+err.Error())
		return
	}

	response.Success(c, config.ToDTO())
}

// UpdateWorkflowConfigAPI 更新工作流配置
func UpdateWorkflowConfigAPI(c *gin.Context) {
	id := c.Param("id")
	configID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "配置ID格式错误")
		return
	}

	var updates WorkflowConfig
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	if err := UpdateWorkflowConfig(uint(configID), &updates); err != nil {
		response.Error(c, http.StatusInternalServerError, "更新配置失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteWorkflowConfigAPI 删除工作流配置
func DeleteWorkflowConfigAPI(c *gin.Context) {
	id := c.Param("id")
	configID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "配置ID格式错误")
		return
	}

	if err := DeleteWorkflowConfig(uint(configID)); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除配置失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// SyncWorkflowConfigsAPI 从 n8n 同步工作流配置
func SyncWorkflowConfigsAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	// 获取所有工作流
	workflowList, err := service.GetWorkflows()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取工作流列表失败: "+err.Error())
		return
	}

	synced := 0
	for _, workflow := range workflowList.Data {
		// 检查是否已存在配置
		_, err := GetWorkflowConfigByWorkflowID(workflow.ID)
		if err == nil {
			// 已存在，跳过
			continue
		}

		// 创建默认配置
		config := &WorkflowConfig{
			WorkflowID:   workflow.ID,
			WorkflowName: workflow.Name,
			Type:         "simple", // 默认类型
			Description:  workflow.Name,
			ExecuteLabel: "执行工作流",
			ResultType:   "auto",
			Fields:       ConfigFields{},
			Enabled:      workflow.Active,
		}

		if err := CreateWorkflowConfig(config); err == nil {
			synced++
		}
	}

	response.Success(c, map[string]interface{}{
		"synced": synced,
		"total":  len(workflowList.Data),
	})
}
