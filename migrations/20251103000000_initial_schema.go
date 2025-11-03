package migrations

import (
	"time"

	"gorm.io/gorm"
)

// Initial20251103000000 初始数据库结构
type Initial20251103000000 struct{}

func (m *Initial20251103000000) Version() string {
	return "20251103000000"
}

func (m *Initial20251103000000) Name() string {
	return "initial_schema"
}

// 用户表
type User struct {
	ID        uint           `gorm:"primarykey"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名"`
	Password  string         `gorm:"type:varchar(255);not null;comment:密码（加密）"`
	Nickname  string         `gorm:"type:varchar(50);comment:昵称"`
	Email     string         `gorm:"type:varchar(100);comment:邮箱"`
	Phone     string         `gorm:"type:varchar(20);comment:手机号"`
	Status    int8           `gorm:"type:tinyint;default:1;comment:状态：0禁用 1启用"`
	CreatedAt time.Time      `gorm:"comment:创建时间"`
	UpdatedAt time.Time      `gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// 用户Token表
type UserToken struct {
	ID           uint      `gorm:"primarykey"`
	UserID       uint      `gorm:"index;not null;comment:用户ID"`
	Token        string    `gorm:"type:varchar(500);uniqueIndex;not null;comment:token值"`
	RefreshToken string    `gorm:"type:varchar(500);comment:刷新token"`
	ExpiresAt    time.Time `gorm:"comment:过期时间"`
	CreatedAt    time.Time `gorm:"comment:创建时间"`
	UpdatedAt    time.Time `gorm:"comment:更新时间"`
}

// 证书监控表
type CertificateMonitor struct {
	ID           uint           `gorm:"primarykey"`
	URL          string         `gorm:"type:varchar(255);not null;comment:监控的URL"`
	Domain       string         `gorm:"type:varchar(255);comment:域名"`
	CustomerName string         `gorm:"type:varchar(100);comment:客户名称"`
	Issuer       string         `gorm:"type:varchar(255);comment:证书颁发者"`
	Subject      string         `gorm:"type:varchar(255);comment:证书主题"`
	Organization string         `gorm:"type:varchar(255);comment:组织"`
	NotBefore    time.Time      `gorm:"comment:证书生效时间"`
	NotAfter     time.Time      `gorm:"comment:证书过期时间"`
	DaysLeft     int            `gorm:"comment:剩余天数"`
	IsValid      bool           `gorm:"comment:是否有效"`
	Status       int            `gorm:"type:tinyint;default:0;comment:状态：0未知 1正常 2警告 3过期"`
	ErrorMsg     string         `gorm:"type:text;comment:错误信息"`
	Remarks      string         `gorm:"type:text;comment:备注"`
	LastCheckAt  time.Time      `gorm:"comment:最后检查时间"`
	CreatedAt    time.Time      `gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `gorm:"comment:更新时间"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

// 密码账户表
type PasswordAccount struct {
	ID              uint           `gorm:"primarykey"`
	UserID          uint           `gorm:"index;not null;comment:用户ID"`
	Title           string         `gorm:"type:varchar(255);not null;comment:标题"`
	URL             string         `gorm:"type:varchar(500);comment:网站地址"`
	Username        string         `gorm:"type:varchar(255);comment:用户名"`
	Password        string         `gorm:"type:text;comment:加密后的密码"`
	EncryptedKey    string         `gorm:"type:varchar(255);comment:加密密钥"`
	Category        string         `gorm:"type:varchar(50);comment:分类"`
	Notes           string         `gorm:"type:text;comment:备注"`
	IsFav           bool           `gorm:"default:false;comment:是否收藏"`
	EncryptionAlgo  string         `gorm:"type:varchar(50);default:'AES-256-GCM';comment:加密算法"`
	CreatedAt       time.Time      `gorm:"comment:创建时间"`
	UpdatedAt       time.Time      `gorm:"comment:更新时间"`
	DeletedAt       gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

func (m *Initial20251103000000) Up(db *gorm.DB) error {
	// 创建所有表
	tables := []interface{}{
		&User{},
		&UserToken{},
		&CertificateMonitor{},
		&PasswordAccount{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return err
		}
	}

	// 创建默认管理员账户（密码：admin123）
	// 密码使用 bcrypt 加密
	defaultUser := &User{
		Username: "admin",
		Password: "$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE/OUisGQlTi3m", // admin123
		Nickname: "管理员",
		Status:   1,
	}

	// 检查是否已存在管理员
	var count int64
	if err := db.Model(&User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		if err := db.Create(defaultUser).Error; err != nil {
			return err
		}
	}

	return nil
}

func (m *Initial20251103000000) Down(db *gorm.DB) error {
	// 删除所有表（注意顺序，先删除有外键关系的表）
	tables := []interface{}{
		&PasswordAccount{},
		&CertificateMonitor{},
		&UserToken{},
		&User{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}

