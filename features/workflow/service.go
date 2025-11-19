package workflow

import (
	"bytes"
	"casper_go/config"
	"casper_go/core/logger"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Service n8n 工作流服务
type Service struct {
	apiURL string
	apiKey string
	client *http.Client
}

// NewService 创建工作流服务实例
func NewService(cfg *config.N8NConfig) *Service {
	return &Service{
		apiURL: cfg.APIURL,
		apiKey: cfg.APIKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest 发送 HTTP 请求到 n8n API
func (s *Service) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := s.apiURL + "/api/v1" + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("X-N8N-API-KEY", s.apiKey)
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

// parseResponse 解析响应
func parseResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var n8nErr N8NError
		if err := json.Unmarshal(body, &n8nErr); err != nil {
			return fmt.Errorf("n8n API 错误 (状态码: %d): %s", resp.StatusCode, string(body))
		}
		return &n8nErr
	}

	// 解析成功响应
	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("解析响应数据失败: %w", err)
		}
	}

	return nil
}

// GetWorkflows 获取所有工作流
func (s *Service) GetWorkflows() (*WorkflowListResponse, error) {
	resp, err := s.makeRequest("GET", "/workflows", nil)
	if err != nil {
		logger.Log.Errorf("获取工作流列表失败: %v", err)
		return nil, err
	}

	var result WorkflowListResponse
	if err := parseResponse(resp, &result); err != nil {
		logger.Log.Errorf("解析工作流列表失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// GetWorkflow 获取单个工作流详情
func (s *Service) GetWorkflow(workflowID string) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s", workflowID)
	resp, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		logger.Log.Errorf("获取工作流详情失败: %v", err)
		return nil, err
	}

	var result Workflow
	if err := parseResponse(resp, &result); err != nil {
		logger.Log.Errorf("解析工作流详情失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// ExecuteWorkflow 执行工作流
func (s *Service) ExecuteWorkflow(workflowID string, data map[string]interface{}) (*ExecutionData, error) {
	// 首先尝试获取工作流信息，检查是否有 Webhook
	workflow, err := s.GetWorkflow(workflowID)
	if err != nil {
		logger.Log.Errorf("获取工作流信息失败: %v", err)
		return nil, err
	}

	// 检查工作流是否使用 Webhook 触发器
	webhookURL := s.findWebhookURL(workflow)
	if webhookURL != "" {
		logger.Log.Infof("检测到 Webhook 触发器，使用 Webhook 执行: %s", webhookURL)
		return s.executeViaWebhook(webhookURL, data)
	}

	// 否则尝试使用 API 执行端点
	endpoint := fmt.Sprintf("/workflows/%s/run", workflowID)
	reqData := data
	if reqData == nil {
		reqData = make(map[string]interface{})
	}

	logger.Log.Infof("尝试执行工作流 %s，使用端点: %s", workflowID, endpoint)

	resp, err := s.makeRequest("POST", endpoint, reqData)
	if err != nil {
		logger.Log.Errorf("执行工作流失败: %v", err)
		return nil, err
	}

	var result ExecutionData
	if err := parseResponse(resp, &result); err != nil {
		logger.Log.Errorf("解析执行结果失败: %v", err)
		return nil, err
	}

	logger.Log.Infof("工作流 %s 执行成功，执行ID: %s", workflowID, result.ID)
	return &result, nil
}

// findWebhookURL 从工作流中查找 Webhook URL
func (s *Service) findWebhookURL(workflow *Workflow) string {
	for _, node := range workflow.Nodes {
		if node.Type == "n8n-nodes-base.webhook" {
			// 构建 Webhook URL
			if path, ok := node.Parameters["path"].(string); ok {
				// 优先使用生产 Webhook URL（永久有效）
				if workflow.Active {
					logger.Log.Infof("工作流已激活，使用生产 Webhook URL")
					return fmt.Sprintf("%s/webhook/%s", s.apiURL, path)
				}
				// 工作流未激活时使用测试 URL（注意：测试 URL 只能用一次）
				if node.WebhookID != "" {
					logger.Log.Warnf("工作流未激活，使用测试 Webhook URL（只能使用一次）")
					return fmt.Sprintf("%s/webhook-test/%s", s.apiURL, node.WebhookID)
				}
			}
		}
	}
	return ""
}

// executeViaWebhook 通过 Webhook 执行工作流
func (s *Service) executeViaWebhook(webhookURL string, data map[string]interface{}) (*ExecutionData, error) {
	// 先尝试 GET 请求（最常见）
	method := "GET"
	var reqBody io.Reader
	
	// 如果有参数，则使用 POST
	if data != nil && len(data) > 0 {
		method = "POST"
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("序列化请求数据失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	logger.Log.Infof("使用 %s 方法调用 Webhook: %s", method, webhookURL)

	req, err := http.NewRequest(method, webhookURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建 Webhook 请求失败: %w", err)
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Webhook 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 Webhook 响应失败: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		logger.Log.Errorf("Webhook 返回错误状态码 %d: %s", resp.StatusCode, string(body))
		
		// 如果是 404 且提示方法不对，尝试另一个方法
		if resp.StatusCode == 404 && method == "POST" {
			logger.Log.Infof("POST 失败，尝试使用 GET 请求")
			return s.executeViaWebhookWithMethod(webhookURL, "GET")
		}
	}

	// Webhook 直接返回工作流结果，需要包装成 ExecutionData 格式
	result := &ExecutionData{
		ID:         fmt.Sprintf("webhook-%d", time.Now().Unix()),
		Finished:   true,
		Mode:       "webhook",
		StartedAt:  time.Now(),
		StoppedAt:  time.Now(),
		WorkflowID: "",
		Status:     "success",
		Data: map[string]interface{}{
			"result": string(body),
		},
	}

	logger.Log.Infof("Webhook 执行成功: %s", webhookURL)
	return result, nil
}

// executeViaWebhookWithMethod 使用指定的 HTTP 方法执行 Webhook
func (s *Service) executeViaWebhookWithMethod(webhookURL, method string) (*ExecutionData, error) {
	req, err := http.NewRequest(method, webhookURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 Webhook 请求失败: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Webhook 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 Webhook 响应失败: %w", err)
	}

	result := &ExecutionData{
		ID:         fmt.Sprintf("webhook-%d", time.Now().Unix()),
		Finished:   true,
		Mode:       "webhook",
		StartedAt:  time.Now(),
		StoppedAt:  time.Now(),
		WorkflowID: "",
		Status:     "success",
		Data: map[string]interface{}{
			"result": string(body),
		},
	}

	return result, nil
}

// GetExecutions 获取执行历史记录
func (s *Service) GetExecutions(limit int) (*ExecutionListResponse, error) {
	endpoint := fmt.Sprintf("/executions?limit=%d", limit)
	resp, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		logger.Log.Errorf("获取执行历史失败: %v", err)
		return nil, err
	}

	var result ExecutionListResponse
	if err := parseResponse(resp, &result); err != nil {
		logger.Log.Errorf("解析执行历史失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// GetExecution 获取单个执行记录详情
func (s *Service) GetExecution(executionID string) (*ExecutionData, error) {
	endpoint := fmt.Sprintf("/executions/%s", executionID)
	resp, err := s.makeRequest("GET", endpoint, nil)
	if err != nil {
		logger.Log.Errorf("获取执行详情失败: %v", err)
		return nil, err
	}

	var result ExecutionData
	if err := parseResponse(resp, &result); err != nil {
		logger.Log.Errorf("解析执行详情失败: %v", err)
		return nil, err
	}

	return &result, nil
}

// DeleteExecution 删除执行记录
func (s *Service) DeleteExecution(executionID string) error {
	endpoint := fmt.Sprintf("/executions/%s", executionID)
	resp, err := s.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		logger.Log.Errorf("删除执行记录失败: %v", err)
		return err
	}

	if err := parseResponse(resp, nil); err != nil {
		logger.Log.Errorf("删除执行记录响应解析失败: %v", err)
		return err
	}

	logger.Log.Infof("执行记录 %s 已删除", executionID)
	return nil
}

// ActivateWorkflow 激活工作流
func (s *Service) ActivateWorkflow(workflowID string) error {
	endpoint := fmt.Sprintf("/workflows/%s", workflowID)

	// 使用 PATCH 方法更新工作流的 active 状态
	body := map[string]interface{}{
		"active": true,
	}

	resp, err := s.makeRequest("PATCH", endpoint, body)
	if err != nil {
		logger.Log.Errorf("激活工作流失败: %v", err)
		return err
	}

	if err := parseResponse(resp, nil); err != nil {
		logger.Log.Errorf("激活工作流响应解析失败: %v", err)
		return err
	}

	logger.Log.Infof("工作流 %s 已激活", workflowID)
	return nil
}

// DeactivateWorkflow 停用工作流
func (s *Service) DeactivateWorkflow(workflowID string) error {
	endpoint := fmt.Sprintf("/workflows/%s", workflowID)

	// 使用 PATCH 方法更新工作流的 active 状态
	body := map[string]interface{}{
		"active": false,
	}

	resp, err := s.makeRequest("PATCH", endpoint, body)
	if err != nil {
		logger.Log.Errorf("停用工作流失败: %v", err)
		return err
	}

	if err := parseResponse(resp, nil); err != nil {
		logger.Log.Errorf("停用工作流响应解析失败: %v", err)
		return err
	}

	logger.Log.Infof("工作流 %s 已停用", workflowID)
	return nil
}

// CheckHealth 检查 n8n 服务健康状态
func (s *Service) CheckHealth() error {
	// n8n 使用 /healthz 作为健康检查端点
	resp, err := s.makeRequest("GET", "/healthz", nil)
	if err != nil {
		return fmt.Errorf("n8n 服务不可用: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("n8n 服务状态异常: %d", resp.StatusCode)
	}

	return nil
}

// ExecuteWorkflowWithFiles 执行包含文件上传的工作流
func (s *Service) ExecuteWorkflowWithFiles(workflowID string, c *gin.Context) (*ExecutionData, error) {
	logger.Log.Infof("开始执行包含文件的工作流: %s", workflowID)

	// 获取工作流信息以找到 Webhook URL
	workflow, err := s.GetWorkflow(workflowID)
	if err != nil {
		return nil, fmt.Errorf("获取工作流信息失败: %w", err)
	}

	// 查找 Webhook URL
	webhookURL := s.findWebhookURL(workflow)
	if webhookURL == "" {
		return nil, fmt.Errorf("未找到 Webhook 触发器，无法执行文件上传")
	}

	logger.Log.Infof("使用 Webhook URL: %s", webhookURL)

	// 创建一个新的 multipart writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 解析表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		return nil, fmt.Errorf("解析表单失败: %w", err)
	}

	// 复制所有文件
	for key, files := range c.Request.MultipartForm.File {
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				return nil, fmt.Errorf("打开文件失败: %w", err)
			}
			defer file.Close()

			part, err := writer.CreateFormFile(key, fileHeader.Filename)
			if err != nil {
				return nil, fmt.Errorf("创建表单文件失败: %w", err)
			}

			if _, err := io.Copy(part, file); err != nil {
				return nil, fmt.Errorf("复制文件失败: %w", err)
			}
		}
	}

	// 复制所有表单字段
	for key, values := range c.Request.MultipartForm.Value {
		for _, value := range values {
			writer.WriteField(key, value)
		}
	}

	writer.Close()

	// 创建请求 - 先尝试 POST
	req, err := http.NewRequest("POST", webhookURL, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Webhook 执行失败: %s\n\n提示：文件上传需要 Webhook 配置为 POST 方法。请在 n8n 中修改 Webhook 节点的 HTTP Method 为 POST", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Webhook 执行失败: %s", string(body))
	}

	// 检查响应类型
	contentType := resp.Header.Get("Content-Type")
	logger.Log.Infof("Webhook 响应 Content-Type: %s", contentType)

	// 检查是否是 Excel 文件
	isExcelFile := contentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
		contentType == "application/vnd.ms-excel" ||
		contentType == "application/octet-stream"

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// Webhook 直接返回工作流结果，需要包装成 ExecutionData 格式
	result := &ExecutionData{
		ID:         fmt.Sprintf("webhook-%d", time.Now().Unix()),
		Finished:   true,
		Mode:       "webhook",
		StartedAt:  time.Now(),
		StoppedAt:  time.Now(),
		WorkflowID: workflowID,
		Status:     "success",
		Data: map[string]interface{}{
			"result":      string(body),
			"contentType": contentType,
			"isFile":      isExcelFile,
		},
	}

	logger.Log.Infof("工作流执行成功: %s, 返回类型: %s, 是否文件: %v", workflowID, contentType, isExcelFile)
	return result, nil
}
