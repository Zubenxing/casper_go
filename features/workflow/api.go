package workflow

import (
	"casper_go/config"
	"casper_go/core/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var service *Service

// InitService 初始化工作流服务
func InitService(cfg *config.N8NConfig) {
	service = NewService(cfg)
	// 注意：工作流配置表由 migrate 系统管理，见 migrations/20251111000007_create_workflow_configs_table.go
}

// GetWorkflowsAPI 获取所有工作流
// @Summary 获取工作流列表
// @Description 获取所有 n8n 工作流
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=WorkflowListResponse}
// @Failure 500 {object} response.Response
// @Router /workflows [get]
func GetWorkflowsAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	result, err := service.GetWorkflows()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取工作流列表失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetWorkflowAPI 获取单个工作流详情
// @Summary 获取工作流详情
// @Description 根据 ID 获取工作流详细信息
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "工作流ID"
// @Success 200 {object} response.Response{data=Workflow}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/{id} [get]
func GetWorkflowAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		response.Error(c, http.StatusBadRequest, "工作流ID不能为空")
		return
	}

	result, err := service.GetWorkflow(workflowID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取工作流详情失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// ExecuteWorkflowAPI 执行工作流
// @Summary 执行工作流
// @Description 执行指定的工作流
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "工作流ID"
// @Param data body map[string]interface{} false "执行参数"
// @Success 200 {object} response.Response{data=ExecutionData}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/{id}/execute [post]
func ExecuteWorkflowAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		response.Error(c, http.StatusBadRequest, "工作流ID不能为空")
		return
	}

	// 检查是否是 multipart/form-data（文件上传）
	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		// 处理文件上传
		result, err := service.ExecuteWorkflowWithFiles(workflowID, c)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "执行工作流失败: "+err.Error())
			return
		}
		response.Success(c, result)
		return
	}

	// 解析请求参数（JSON）
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		// 如果没有请求体，使用空 map
		data = make(map[string]interface{})
	}

	result, err := service.ExecuteWorkflow(workflowID, data)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "执行工作流失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetExecutionsAPI 获取执行历史
// @Summary 获取执行历史
// @Description 获取工作流执行历史记录
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "限制返回数量" default(20)
// @Success 200 {object} response.Response{data=ExecutionListResponse}
// @Failure 500 {object} response.Response
// @Router /workflows/executions [get]
func GetExecutionsAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	result, err := service.GetExecutions(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取执行历史失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// GetExecutionAPI 获取执行详情
// @Summary 获取执行详情
// @Description 根据执行ID获取详细信息
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "执行ID"
// @Success 200 {object} response.Response{data=ExecutionData}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/executions/{id} [get]
func GetExecutionAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	executionID := c.Param("id")
	if executionID == "" {
		response.Error(c, http.StatusBadRequest, "执行ID不能为空")
		return
	}

	result, err := service.GetExecution(executionID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取执行详情失败: "+err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteExecutionAPI 删除执行记录
// @Summary 删除执行记录
// @Description 删除指定的执行记录
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "执行ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/executions/{id} [delete]
func DeleteExecutionAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	executionID := c.Param("id")
	if executionID == "" {
		response.Error(c, http.StatusBadRequest, "执行ID不能为空")
		return
	}

	if err := service.DeleteExecution(executionID); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除执行记录失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// ActivateWorkflowAPI 激活工作流
// @Summary 激活工作流
// @Description 激活指定的工作流
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "工作流ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/{id}/activate [post]
func ActivateWorkflowAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		response.Error(c, http.StatusBadRequest, "工作流ID不能为空")
		return
	}

	if err := service.ActivateWorkflow(workflowID); err != nil {
		response.Error(c, http.StatusInternalServerError, "激活工作流失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "工作流已激活"})
}

// DeactivateWorkflowAPI 停用工作流
// @Summary 停用工作流
// @Description 停用指定的工作流
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "工作流ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/{id}/deactivate [post]
func DeactivateWorkflowAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	workflowID := c.Param("id")
	if workflowID == "" {
		response.Error(c, http.StatusBadRequest, "工作流ID不能为空")
		return
	}

	if err := service.DeactivateWorkflow(workflowID); err != nil {
		response.Error(c, http.StatusInternalServerError, "停用工作流失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"message": "工作流已停用"})
}

// CheckHealthAPI 检查 n8n 服务健康状态
// @Summary 检查服务健康状态
// @Description 检查 n8n 服务是否正常运行
// @Tags workflow
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /workflows/health [get]
func CheckHealthAPI(c *gin.Context) {
	if service == nil {
		response.Error(c, http.StatusInternalServerError, "工作流服务未初始化")
		return
	}

	if err := service.CheckHealth(); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "n8n 服务不可用: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"status":  "healthy",
		"message": "n8n 服务运行正常",
	})
}
