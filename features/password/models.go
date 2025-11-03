package password

import (
	"time"

	"gorm.io/gorm"
)

// Account 账户密码模型
type Account struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID   uint   `gorm:"not null;index" json:"user_id"`     // 所属用户ID
	Title    string `gorm:"size:100;not null" json:"title"`    // 网站/应用标题
	URL      string `gorm:"size:500" json:"url"`               // 网站URL
	Username string `gorm:"size:200" json:"username"`          // 用户名/账号
	Password string `gorm:"type:text;not null" json:"-"`       // 密码（加密存储，不返回给前端）
	Category string `gorm:"size:50;index" json:"category"`     // 分类（如：社交、邮箱、工作等）
	Notes    string `gorm:"type:text" json:"notes"`            // 备注
	IsFav    bool   `gorm:"default:false;index" json:"is_fav"` // 是否收藏
}

// TableName 指定表名
func (Account) TableName() string {
	return "password_accounts"
}

// CreateRequest 创建账户请求
type CreateRequest struct {
	Title    string `json:"title" binding:"required"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
	IsFav    bool   `json:"is_fav"`
}

// UpdateRequest 更新账户请求
type UpdateRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"` // 如果为空则不更新密码
	Category string `json:"category"`
	Notes    string `json:"notes"`
	IsFav    *bool  `json:"is_fav"` // 使用指针以区分未设置和false
}

// Response 账户响应（不包含密码）
type Response struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Username  string    `json:"username"`
	Category  string    `json:"category"`
	Notes     string    `json:"notes"`
	IsFav     bool      `json:"is_fav"`
}

// ToResponse 转换为响应格式
func (a *Account) ToResponse() *Response {
	return &Response{
		ID:        a.ID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Title:     a.Title,
		URL:       a.URL,
		Username:  a.Username,
		Category:  a.Category,
		Notes:     a.Notes,
		IsFav:     a.IsFav,
	}
}

// PasswordResponse 密码响应（用于查看密码时）
type PasswordResponse struct {
	Password string `json:"password"`
}

// ListResponse 列表响应
type ListResponse struct {
	Total int64       `json:"total"`
	Items []*Response `json:"items"`
}
