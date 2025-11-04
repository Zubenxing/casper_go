package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 全局配置结构
type Config struct {
	Database     DatabaseConfig     `mapstructure:"database"`
	DatabaseInit DatabaseInitConfig `mapstructure:"database_init"`
	Site         SiteConfig         `mapstructure:"site"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Logger       LoggerConfig       `mapstructure:"logger"`
	API          APIConfig          `mapstructure:"api"`
	Certificate  CertificateConfig  `mapstructure:"certificate"`
	Upload       UploadConfig       `mapstructure:"upload"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	LogLevel        string `mapstructure:"log_level"`
}

// DatabaseInitConfig 数据库初始化配置
type DatabaseInitConfig struct {
	AutoMigrate     bool `mapstructure:"auto_migrate"`
	InitDefaultData bool `mapstructure:"init_default_data"`
}

// SiteConfig 站点配置
type SiteConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
}

// APIConfig API 配置
type APIConfig struct {
	EnableVersion     bool   `mapstructure:"enable_version"`
	CurrentVersion    string `mapstructure:"current_version"`
	RateLimitEnabled  bool   `mapstructure:"rate_limit_enabled"`
	RateLimitRequests int    `mapstructure:"rate_limit_requests"`
	RequestIDEnabled  bool   `mapstructure:"request_id_enabled"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret             string `mapstructure:"secret"`
	ExpireHours        int    `mapstructure:"expire_hours"`
	RefreshExpireHours int    `mapstructure:"refresh_expire_hours"`
	Issuer             string `mapstructure:"issuer"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// CertificateConfig 证书监控配置
type CertificateConfig struct {
	CheckInterval int `mapstructure:"check_interval"`
	WarningDays   int `mapstructure:"warning_days"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	BaseDir       string `mapstructure:"base_dir"`
	WorkIssuesDir string `mapstructure:"work_issues_dir"`
	MaxFileSize   int    `mapstructure:"max_file_size"` // MB
	AllowedTypes  string `mapstructure:"allowed_types"`
}

// GetWorkIssuesPath 获取工作问题图片上传路径
func (c *UploadConfig) GetWorkIssuesPath() string {
	return filepath.Join(c.BaseDir, c.WorkIssuesDir)
}

// GetMaxFileSizeBytes 获取最大文件大小（字节）
func (c *UploadConfig) GetMaxFileSizeBytes() int64 {
	return int64(c.MaxFileSize) * 1024 * 1024
}

var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 加载数据库配置
	dbViper := viper.New()
	dbViper.SetConfigName("database")
	dbViper.SetConfigType("yaml")
	dbViper.AddConfigPath(configPath)

	if err := dbViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取数据库配置失败: %w", err)
	}

	// 加载站点配置
	siteViper := viper.New()
	siteViper.SetConfigName("site_info")
	siteViper.SetConfigType("yaml")
	siteViper.AddConfigPath(configPath)

	if err := siteViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取站点配置失败: %w", err)
	}

	// 合并配置
	cfg := &Config{}

	if err := dbViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析数据库配置失败: %w", err)
	}

	if err := siteViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析站点配置失败: %w", err)
	}

	GlobalConfig = cfg
	return cfg, nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
	)
}

// GetLogPath 获取日志文件的完整路径
func (c *LoggerConfig) GetLogPath() string {
	if filepath.IsAbs(c.Output) {
		return c.Output
	}
	return filepath.Join(".", c.Output)
}
