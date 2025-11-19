package workflow

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// WorkflowConfig 工作流前端配置
type WorkflowConfig struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	WorkflowID     string         `json:"workflowId" gorm:"type:varchar(100);uniqueIndex;not null;comment:n8n工作流ID"`
	WorkflowName   string         `json:"workflowName" gorm:"type:varchar(255);not null;comment:工作流名称"`
	Type           string         `json:"type" gorm:"type:varchar(50);not null;default:simple;comment:工作流类型"`
	Description    string         `json:"description" gorm:"type:text;comment:描述"`
	ExecuteLabel   string         `json:"executeLabel" gorm:"type:varchar(50);default:执行;comment:执行按钮文本"`
	ResultType     string         `json:"resultType" gorm:"type:varchar(50);default:auto;comment:结果类型"`
	Fields         ConfigFields   `json:"fields" gorm:"type:text;comment:字段配置"`
	Enabled        bool           `json:"enabled" gorm:"default:true;comment:是否启用"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// ConfigField 字段配置
type ConfigField struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"` // file, text, textarea, number, date
	Accept      string `json:"accept,omitempty"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
}

// ConfigFields 字段数组类型
type ConfigFields []ConfigField

// Value 实现 driver.Valuer 接口
func (f ConfigFields) Value() (driver.Value, error) {
	if f == nil {
		return "[]", nil
	}
	return json.Marshal(f)
}

// Scan 实现 sql.Scanner 接口
func (f *ConfigFields) Scan(value interface{}) error {
	if value == nil {
		*f = ConfigFields{}
		return nil
	}
	
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	
	return json.Unmarshal(bytes, f)
}

// TableName 指定表名
func (WorkflowConfig) TableName() string {
	return "workflow_configs"
}

// WorkflowConfigDTO 工作流配置 DTO
type WorkflowConfigDTO struct {
	ID           uint         `json:"id"`
	WorkflowID   string       `json:"workflowId"`
	WorkflowName string       `json:"workflowName"`
	Type         string       `json:"type"`
	Description  string       `json:"description"`
	ExecuteLabel string       `json:"executeLabel"`
	ResultType   string       `json:"resultType"`
	Fields       ConfigFields `json:"fields"`
	Enabled      bool         `json:"enabled"`
}

// ToDTO 转换为 DTO
func (w *WorkflowConfig) ToDTO() *WorkflowConfigDTO {
	return &WorkflowConfigDTO{
		ID:           w.ID,
		WorkflowID:   w.WorkflowID,
		WorkflowName: w.WorkflowName,
		Type:         w.Type,
		Description:  w.Description,
		ExecuteLabel: w.ExecuteLabel,
		ResultType:   w.ResultType,
		Fields:       w.Fields,
		Enabled:      w.Enabled,
	}
}
