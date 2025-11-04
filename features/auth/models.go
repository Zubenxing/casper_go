package auth

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null;comment:用户ID(1000-2000:admin, 2000-10000:普通用户)" json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nickname string `gorm:"type:varchar(100)" json:"nickname"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	Role     string `gorm:"type:varchar(20);default:'user'" json:"role"` // admin, user
	Status   int    `gorm:"type:tinyint;default:1" json:"status"`        // 1:启用 0:禁用

	// 注意：Tokens 关联已移除，因为 user_id 现在是独立字段，不再使用外键
	// Tokens []UserToken `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserToken 用户令牌模型
type UserToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	UserID           uint      `gorm:"index;not null;comment:用户ID(对应users.user_id，非外键)" json:"user_id"`
	Token            string    `gorm:"type:varchar(500);uniqueIndex;not null" json:"token"`
	RefreshToken     string    `gorm:"type:varchar(500);uniqueIndex" json:"refresh_token"`
	TokenType        string    `gorm:"type:varchar(20);default:'Bearer'" json:"token_type"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	IPAddress        string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent        string    `gorm:"type:varchar(255)" json:"user_agent"`
	IsValid          bool      `gorm:"default:true" json:"is_valid"`

	// 注意：User 关联已移除，因为 user_id 现在对应 users.user_id，不再使用外键
	// User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (UserToken) TableName() string {
	return "user_tokens"
}

// IsExpired 检查 token 是否过期
func (t *UserToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsRefreshExpired 检查 refresh token 是否过期
func (t *UserToken) IsRefreshExpired() bool {
	return time.Now().After(t.RefreshExpiresAt)
}

// ========== DTO 对象 ==========

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse Token 响应
type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	ExpiresAt    time.Time     `json:"expires_at"`
	User         *UserResponse `json:"user"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		UserID:    u.UserID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
