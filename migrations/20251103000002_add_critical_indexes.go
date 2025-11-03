package migrations

import (
	"gorm.io/gorm"
)

// AddCriticalIndexes20251103000002 添加关键性能索引
type AddCriticalIndexes20251103000002 struct{}

func (m *AddCriticalIndexes20251103000002) Version() string {
	return "20251103000002"
}

func (m *AddCriticalIndexes20251103000002) Name() string {
	return "add_critical_indexes_for_jwt_and_soft_delete"
}

func (m *AddCriticalIndexes20251103000002) Up(db *gorm.DB) error {
	// 关键索引配置
	indexes := []struct {
		table   string
		name    string
		sql     string
		comment string
	}{
		{
			table:   "user_tokens",
			name:    "idx_token_valid_expires",
			sql:     "CREATE INDEX idx_token_valid_expires ON user_tokens(token(255), is_valid, expires_at)",
			comment: "JWT验证核心索引 - 覆盖token验证的所有条件",
		},
		{
			table:   "password_accounts",
			name:    "idx_user_deleted",
			sql:     "CREATE INDEX idx_user_deleted ON password_accounts(user_id, deleted_at)",
			comment: "用户ID+软删除复合索引 - 加速所有用户查询",
		},
		{
			table:   "password_accounts",
			name:    "idx_user_category_deleted",
			sql:     "CREATE INDEX idx_user_category_deleted ON password_accounts(user_id, category, deleted_at)",
			comment: "用户+分类+软删除索引 - 加速分类筛选",
		},
		{
			table:   "password_accounts",
			name:    "idx_user_fav_created_deleted",
			sql:     "CREATE INDEX idx_user_fav_created_deleted ON password_accounts(user_id, is_fav, created_at, deleted_at)",
			comment: "用户+收藏+创建时间+软删除索引 - 加速列表排序",
		},
		{
			table:   "certificate_monitors",
			name:    "idx_deleted_at",
			sql:     "CREATE INDEX idx_deleted_at ON certificate_monitors(deleted_at)",
			comment: "证书监控软删除索引",
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

func (m *AddCriticalIndexes20251103000002) Down(db *gorm.DB) error {
	// 删除索引
	indexes := []struct {
		table string
		name  string
	}{
		{"user_tokens", "idx_token_valid_expires"},
		{"password_accounts", "idx_user_deleted"},
		{"password_accounts", "idx_user_category_deleted"},
		{"password_accounts", "idx_user_fav_created_deleted"},
		{"certificate_monitors", "idx_deleted_at"},
	}

	for _, idx := range indexes {
		sql := "DROP INDEX " + idx.name + " ON " + idx.table
		// 忽略索引不存在的错误
		db.Exec(sql)
	}

	return nil
}
