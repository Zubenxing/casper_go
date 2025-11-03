package migrations

import (
	"gorm.io/gorm"
)

// OptimizePasswordIndexes20251103000003 优化密码管理索引
type OptimizePasswordIndexes20251103000003 struct{}

func (m *OptimizePasswordIndexes20251103000003) Version() string {
	return "20251103000003"
}

func (m *OptimizePasswordIndexes20251103000003) Name() string {
	return "optimize_password_indexes_for_query_performance"
}

func (m *OptimizePasswordIndexes20251103000003) Up(db *gorm.DB) error {
	// 删除一些冗余或低效的索引
	redundantIndexes := []struct {
		table string
		name  string
	}{
		// 这些索引被更全面的复合索引覆盖
		{"password_accounts", "idx_password_accounts_deleted_at"},
		{"password_accounts", "idx_password_accounts_user_id"},
		{"password_accounts", "idx_user_fav"},
		{"password_accounts", "idx_created_at"},
	}

	for _, idx := range redundantIndexes {
		// 先检查索引是否存在
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name).Scan(&count)

		if count > 0 {
			sql := "DROP INDEX " + idx.name + " ON " + idx.table
			if err := db.Exec(sql).Error; err != nil {
				// 忽略索引不存在的错误
				continue
			}
		}
	}

	// 优化后的索引结构说明：
	// 1. idx_user_deleted - 用于基本的用户查询和软删除过滤
	// 2. idx_user_category_deleted - 用于分类筛选
	// 3. idx_user_fav_created_deleted - 用于列表排序（收藏+时间）
	// 这三个索引覆盖了所有查询场景，无需其他冗余索引

	return nil
}

func (m *OptimizePasswordIndexes20251103000003) Down(db *gorm.DB) error {
	// 恢复被删除的索引
	indexes := []struct {
		table string
		name  string
		sql   string
	}{
		{
			table: "password_accounts",
			name:  "idx_password_accounts_deleted_at",
			sql:   "CREATE INDEX idx_password_accounts_deleted_at ON password_accounts(deleted_at)",
		},
		{
			table: "password_accounts",
			name:  "idx_password_accounts_user_id",
			sql:   "CREATE INDEX idx_password_accounts_user_id ON password_accounts(user_id)",
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
	}

	for _, idx := range indexes {
		var count int64
		db.Raw("SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name).Scan(&count)

		if count == 0 {
			db.Exec(idx.sql)
		}
	}

	return nil
}
