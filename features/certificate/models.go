package certificate

import (
	"time"

	"gorm.io/gorm"
)

// Monitor 证书监控模型
type Monitor struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	URL          string    `gorm:"type:varchar(500);uniqueIndex;not null" json:"url"`
	Domain       string    `gorm:"type:varchar(255)" json:"domain"`
	CustomerName string    `gorm:"type:varchar(255)" json:"customer_name"` // 客户名
	Issuer       string    `gorm:"type:varchar(255)" json:"issuer"`
	Subject      string    `gorm:"type:varchar(255)" json:"subject"`
	Organization string    `gorm:"type:varchar(255)" json:"organization"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	DaysLeft     int       `json:"days_left"`
	IsValid      bool      `gorm:"default:true" json:"is_valid"`
	ErrorMsg     string    `gorm:"type:text" json:"error_msg"`
	Remark       string    `gorm:"type:text" json:"remark"` // 备注
	LastCheckAt  time.Time `json:"last_check_at"`
	Status       int       `gorm:"type:tinyint;default:1" json:"status"` // 1:正常 2:警告 3:过期 0:禁用
}

// TableName 指定表名
func (Monitor) TableName() string {
	return "certificate_monitors"
}

// ========== DTO 对象 ==========

// AddRequest 添加证书监控请求
type AddRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// Response 证书响应
type Response struct {
	ID           uint      `json:"id"`
	URL          string    `json:"url"`
	Domain       string    `json:"domain"`
	CustomerName string    `json:"customer_name"`
	Issuer       string    `json:"issuer"`
	Subject      string    `json:"subject"`
	Organization string    `json:"organization"`
	NotBefore    time.Time `json:"not_before"`
	NotAfter     time.Time `json:"not_after"`
	DaysLeft     int       `json:"days_left"`
	IsValid      bool      `json:"is_valid"`
	Status       int       `json:"status"`
	StatusText   string    `json:"status_text"`
	ErrorMsg     string    `json:"error_msg,omitempty"`
	Remark       string    `json:"remark,omitempty"`
	LastCheckAt  time.Time `json:"last_check_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateRequest 更新证书信息请求
type UpdateRequest struct {
	CustomerName string `json:"customer_name"`
	Remark       string `json:"remark"`
}

// ToResponse 转换为响应格式
func (m *Monitor) ToResponse() *Response {
	statusText := ""
	switch m.Status {
	case 1:
		statusText = "正常"
	case 2:
		statusText = "警告"
	case 3:
		statusText = "过期"
	case 0:
		statusText = "禁用"
	}

	return &Response{
		ID:           m.ID,
		URL:          m.URL,
		Domain:       m.Domain,
		CustomerName: m.CustomerName,
		Issuer:       m.Issuer,
		Subject:      m.Subject,
		Organization: m.Organization,
		NotBefore:    m.NotBefore,
		NotAfter:     m.NotAfter,
		DaysLeft:     m.DaysLeft,
		IsValid:      m.IsValid,
		Status:       m.Status,
		StatusText:   statusText,
		ErrorMsg:     m.ErrorMsg,
		Remark:       m.Remark,
		LastCheckAt:  m.LastCheckAt,
		CreatedAt:    m.CreatedAt,
	}
}
