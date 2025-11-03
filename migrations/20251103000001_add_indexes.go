package migrations

import (
	"gorm.io/gorm"
)

// AddIndexes20251103000001 添加索引优化查询性能
type AddIndexes20251103000001 struct{}

func (m *AddIndexes20251103000001) Version() string {
	return "20251103000001"
}

func (m *AddIndexes20251103000001) Name() string {
	return "add_indexes_for_performance"
}

func (m *AddIndexes20251103000001) Up(db *gorm.DB) error {
	// 证书监控表索引
	indexes := []struct {
		table string
		name  string
		sql   string
	}{
		{
			table: "certificate_monitors",
			name:  "idx_status_days_left",
			sql:   "CREATE INDEX idx_status_days_left ON certificate_monitors(status, days_left)",
		},
		{
			table: "certificate_monitors",
			name:  "idx_last_check_at",
			sql:   "CREATE INDEX idx_last_check_at ON certificate_monitors(last_check_at)",
		},
		{
			table: "certificate_monitors",
			name:  "idx_url",
			sql:   "CREATE INDEX idx_url ON certificate_monitors(url)",
		},
		{
			table: "certificate_monitors",
			name:  "idx_domain",
			sql:   "CREATE INDEX idx_domain ON certificate_monitors(domain)",
		},
		// 密码账户表索引
		{
			table: "password_accounts",
			name:  "idx_user_category",
			sql:   "CREATE INDEX idx_user_category ON password_accounts(user_id, category)",
		},
		{
			table: "password_accounts",
			name:  "idx_user_fav",
			sql:   "CREATE INDEX idx_user_fav ON password_accounts(user_id, is_fav)",
		},
		{
			table: "password_accounts",
			name:  "idx_created_at",
			sql:   "CREATE INDEX idx_created_at ON password_accounts(created_at)",
		},
		// 用户Token表索引
		{
			table: "user_tokens",
			name:  "idx_expires_at",
			sql:   "CREATE INDEX idx_expires_at ON user_tokens(expires_at)",
		},
	}

	for _, idx := range indexes {
		// 检查索引是否已存在
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name).Scan(&count)

		if count == 0 {
			if err := db.Exec(idx.sql).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *AddIndexes20251103000001) Down(db *gorm.DB) error {
	// 删除索引
	indexes := []struct {
		table string
		name  string
	}{
		{"certificate_monitors", "idx_status_days_left"},
		{"certificate_monitors", "idx_last_check_at"},
		{"certificate_monitors", "idx_url"},
		{"certificate_monitors", "idx_domain"},
		{"password_accounts", "idx_user_category"},
		{"password_accounts", "idx_user_fav"},
		{"password_accounts", "idx_created_at"},
		{"user_tokens", "idx_expires_at"},
	}

	for _, idx := range indexes {
		sql := "DROP INDEX " + idx.name + " ON " + idx.table
		// 忽略索引不存在的错误
		db.Exec(sql)
	}

	return nil
}

