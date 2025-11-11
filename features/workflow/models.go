package workflow

import "time"

// Workflow n8n 工作流模型
type Workflow struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Nodes     []Node    `json:"nodes"`
	WebhookURL string    `json:"webhookUrl,omitempty"` // Webhook URL (如果有)
}

// Tag 工作流标签
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Node 工作流节点
type Node struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	TypeVersion float64                `json:"typeVersion"`
	Position    []float64              `json:"position"`
	Parameters  map[string]interface{} `json:"parameters"`
	WebhookID   string                 `json:"webhookId,omitempty"` // Webhook ID
}

// WorkflowExecution 工作流执行记录
type WorkflowExecution struct {
	ID             string    `json:"id"`
	Finished       bool      `json:"finished"`
	Mode           string    `json:"mode"`
	StartedAt      time.Time `json:"startedAt"`
	StoppedAt      time.Time `json:"stoppedAt"`
	WorkflowID     string    `json:"workflowId"`
	WorkflowName   string    `json:"workflowName"`
	Status         string    `json:"status"` // success, error, waiting
	RetryOf        string    `json:"retryOf"`
	RetrySuccessId string    `json:"retrySuccessId"`
}

// ExecutionData 执行详情数据
type ExecutionData struct {
	ID         string                 `json:"id"`
	Data       map[string]interface{} `json:"data"`
	Finished   bool                   `json:"finished"`
	Mode       string                 `json:"mode"`
	StartedAt  time.Time              `json:"startedAt"`
	StoppedAt  time.Time              `json:"stoppedAt"`
	WorkflowID string                 `json:"workflowId"`
	Status     string                 `json:"status"`
}

// ExecuteWorkflowRequest 执行工作流请求
type ExecuteWorkflowRequest struct {
	WorkflowID string                 `json:"workflowId" binding:"required"`
	Data       map[string]interface{} `json:"data"`
}

// WorkflowListResponse 工作流列表响应
type WorkflowListResponse struct {
	Data  []Workflow `json:"data"`
	Count int        `json:"count"`
}

// ExecutionListResponse 执行记录列表响应
type ExecutionListResponse struct {
	Data  []WorkflowExecution `json:"data"`
	Count int                 `json:"count"`
}

// N8NError n8n API 错误响应
type N8NError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *N8NError) Error() string {
	return e.Message
}
