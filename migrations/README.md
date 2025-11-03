# 数据库迁移系统

## 概述

本项目使用自定义的数据库迁移系统，类似于 Laravel 的 migration 机制，用于版本化管理数据库结构变更。

## 迁移文件命名规范

```
YYYYMMDDHHMMSS_description.go
```

例如：
- `20251103000000_initial_schema.go` - 初始数据库结构
- `20251103000001_add_indexes.go` - 添加索引
- `20251103120000_add_user_avatar.go` - 添加用户头像字段

## 使用方法

### 1. 执行所有待执行的迁移

```bash
# 方式一：使用 go run
go run cmd/migrate/main.go -action=up

# 方式二：编译后执行
go build -o migrate cmd/migrate/main.go
./migrate -action=up

# 指定配置文件路径
go run cmd/migrate/main.go -action=up -config=config
```

### 2. 回滚最后一次迁移

```bash
go run cmd/migrate/main.go -action=down
```

### 3. 查看迁移状态

```bash
go run cmd/migrate/main.go -action=status
```

输出示例：
```
迁移状态:
----------------------------------------
[已执行] 20251103000000 - initial_schema
[未执行] 20251103000001 - add_indexes_for_performance
----------------------------------------
```

## 创建新的迁移

### 1. 创建迁移文件

在 `migrations/` 目录下创建新文件，命名格式：`YYYYMMDDHHMMSS_description.go`

### 2. 实现迁移接口

```go
package migrations

import (
	"gorm.io/gorm"
)

// YourMigration20251103120000 添加新功能
type YourMigration20251103120000 struct{}

func (m *YourMigration20251103120000) Version() string {
	return "20251103120000"
}

func (m *YourMigration20251103120000) Name() string {
	return "add_new_feature"
}

func (m *YourMigration20251103120000) Up(db *gorm.DB) error {
	// 执行数据库变更
	// 例如：添加字段、创建表、插入数据等
	return db.Exec("ALTER TABLE users ADD COLUMN avatar VARCHAR(255)").Error
}

func (m *YourMigration20251103120000) Down(db *gorm.DB) error {
	// 回滚数据库变更
	return db.Exec("ALTER TABLE users DROP COLUMN avatar").Error
}
```

### 3. 注册迁移

在 `migrations/migrations.go` 中注册新迁移：

```go
func GetAllMigrations() []database.Migration {
	return []database.Migration{
		&Initial20251103000000{},
		&AddIndexes20251103000001{},
		&YourMigration20251103120000{}, // 添加这里
	}
}
```

## 迁移最佳实践

### 1. 版本号生成

使用当前时间生成版本号：

```bash
# Linux/Mac
date +%Y%m%d%H%M%S

# PowerShell
Get-Date -Format "yyyyMMddHHmmss"
```

### 2. 迁移内容建议

- ✅ **DO**: 每个迁移只做一件事
- ✅ **DO**: 提供完整的 `Down()` 回滚逻辑
- ✅ **DO**: 在事务中执行（系统自动处理）
- ✅ **DO**: 添加注释说明变更目的
- ❌ **DON'T**: 修改已执行的迁移文件
- ❌ **DON'T**: 删除已执行的迁移文件
- ❌ **DON'T**: 在迁移中使用业务逻辑模型（可能会变化）

### 3. 常见迁移示例

#### 添加新字段

```go
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE users 
		ADD COLUMN avatar VARCHAR(255) COMMENT '头像' AFTER nickname
	`).Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users DROP COLUMN avatar").Error
}
```

#### 创建新表

```go
type NewTable struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
}

func (m *Migration) Up(db *gorm.DB) error {
	return db.AutoMigrate(&NewTable{})
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&NewTable{})
}
```

#### 添加索引

```go
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec("CREATE INDEX idx_email ON users(email)").Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("DROP INDEX idx_email ON users").Error
}
```

#### 数据迁移

```go
func (m *Migration) Up(db *gorm.DB) error {
	// 批量更新数据
	return db.Exec(`
		UPDATE users 
		SET status = 1 
		WHERE status IS NULL
	`).Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec(`
		UPDATE users 
		SET status = NULL 
		WHERE status = 1
	`).Error
}
```

## 迁移历史

系统会在数据库中创建 `migration_histories` 表来记录已执行的迁移：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| version | varchar(20) | 迁移版本号 |
| name | varchar(255) | 迁移名称 |
| created_at | timestamp | 执行时间 |

## 注意事项

1. **生产环境**：执行迁移前请先备份数据库
2. **团队协作**：迁移文件一旦提交到版本控制，不要修改
3. **回滚**：`down` 操作只回滚最后一次迁移，谨慎使用
4. **索引**：大表添加索引可能需要较长时间
5. **事务**：所有迁移在事务中执行，失败会自动回滚

## 与现有代码集成

如果你已经在 `main.go` 中使用了 `AutoMigrate`，可以逐步迁移：

1. 保留现有的 `AutoMigrate` 代码（兼容性）
2. 运行 `migrate up` 创建迁移历史表
3. 逐步将表结构变更移到迁移文件中
4. 最终移除 `AutoMigrate` 代码

## 快捷脚本

### Windows (PowerShell)

创建 `migrate.ps1`:

```powershell
param(
    [string]$Action = "up"
)

go run cmd/migrate/main.go -action=$Action
```

使用：
```powershell
.\migrate.ps1 up
.\migrate.ps1 down
.\migrate.ps1 status
```

### Linux/Mac (Bash)

创建 `migrate.sh`:

```bash
#!/bin/bash
ACTION=${1:-up}
go run cmd/migrate/main.go -action=$ACTION
```

使用：
```bash
./migrate.sh up
./migrate.sh down
./migrate.sh status
```

## 故障排除

### 迁移执行失败

如果迁移执行失败，系统会自动回滚该迁移，不会记录到历史表中。可以：

1. 查看错误日志
2. 修复迁移代码
3. 重新执行 `migrate up`

### 迁移历史不一致

如果手动修改了数据库结构，可以手动调整 `migration_histories` 表：

```sql
-- 查看迁移历史
SELECT * FROM migration_histories ORDER BY created_at DESC;

-- 手动标记迁移已执行
INSERT INTO migration_histories (version, name, created_at) 
VALUES ('20251103000001', 'add_indexes_for_performance', NOW());

-- 删除迁移记录
DELETE FROM migration_histories WHERE version = '20251103000001';
```

