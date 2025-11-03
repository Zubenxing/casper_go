package migrations

import (
	"casper_go/core/database"
)

// GetAllMigrations 获取所有迁移
func GetAllMigrations() []database.Migration {
	return []database.Migration{
		&Initial20251103000000{},
		&AddIndexes20251103000001{},
		&AddCriticalIndexes20251103000002{},
		&OptimizePasswordIndexes20251103000003{},
		// 在这里添加新的迁移...
	}
}
