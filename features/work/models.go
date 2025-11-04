package work

import (
	"time"

	"gorm.io/gorm"
)

// Work 工作记录
type Work struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null;index:idx_work_user_status_created,priority:1;comment:用户ID" json:"user_id"`
	Title       string         `gorm:"type:varchar(200);not null;comment:工作标题" json:"title"`
	Description string         `gorm:"type:text;comment:工作描述" json:"description"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending';index:idx_work_user_status_created,priority:2;comment:状态:pending/in_progress/completed/cancelled" json:"status"`
	Priority    string         `gorm:"type:varchar(20);not null;default:'medium';comment:优先级:low/medium/high/urgent" json:"priority"`
	StartDate   *time.Time     `gorm:"comment:开始日期" json:"start_date"`
	DueDate     *time.Time     `gorm:"comment:截止日期" json:"due_date"`
	CompletedAt *time.Time     `gorm:"comment:完成时间" json:"completed_at"`
	Tags        string         `gorm:"type:varchar(500);comment:标签(逗号分隔)" json:"tags"`
}

// TableName 指定表名
func (Work) TableName() string {
	return "works"
}

// WorkIssue 工作问题记录
type WorkIssue struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"not null;index:idx_issue_user_status_created,priority:1;comment:用户ID" json:"user_id"`
	WorkID      *uint          `gorm:"index;comment:关联的工作ID(可选)" json:"work_id,omitempty"`
	Title       string         `gorm:"type:varchar(200);not null;comment:问题标题" json:"title"`
	Description string         `gorm:"type:text;comment:问题描述" json:"description"`
	Images      string         `gorm:"type:json;comment:问题截图(最多3张)" json:"images,omitempty"`
	Content     string         `gorm:"type:text;comment:问题详细内容(富文本)" json:"content,omitempty"`
	Status      string         `gorm:"type:varchar(20);not null;default:'open';index:idx_issue_user_status_created,priority:2;comment:状态:open/in_progress/resolved/closed" json:"status"`
	Severity    string         `gorm:"type:varchar(20);not null;default:'medium';comment:严重程度:low/medium/high/critical" json:"severity"`
	Solution    string         `gorm:"type:text;comment:解决方案" json:"solution"`
	ResolvedAt  *time.Time     `gorm:"comment:解决时间" json:"resolved_at"`
	Tags        string         `gorm:"type:varchar(500);comment:标签(逗号分隔)" json:"tags"`
}

// TableName 指定表名
func (WorkIssue) TableName() string {
	return "work_issues"
}

// 状态常量
const (
	// Work 状态
	WorkStatusPending    = "pending"
	WorkStatusInProgress = "in_progress"
	WorkStatusCompleted  = "completed"
	WorkStatusCancelled  = "cancelled"

	// 优先级
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"

	// Issue 状态
	IssueStatusOpen       = "open"
	IssueStatusInProgress = "in_progress"
	IssueStatusResolved   = "resolved"
	IssueStatusClosed     = "closed"

	// 严重程度
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"
)
