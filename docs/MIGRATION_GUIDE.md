# 数据库迁移使用指南

## 概述

我们实现了类似 Laravel/PHP 的数据库迁移系统，与你之前使用的 migration 表结构完全兼容。

## 迁移表结构

```sql
CREATE TABLE `migrations` (
    `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Primary ID',
    `version` VARCHAR(20) NOT NULL COMMENT '迁移版本号',
    `migration` VARCHAR(255) NOT NULL COMMENT '迁移名称',
    `batch` INT(10) NOT NULL DEFAULT 0 COMMENT '批次号',
    `created_at` DATETIME(3) NULL COMMENT '执行时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `idx_migrations_version` (`version`)
)
```

## 核心功能

### ✅ Batch 批次管理

- **批次号自动递增**：每次执行 `migrate up` 时，所有新迁移都会分配到同一个批次号
- **批量回滚**：`migrate down` 会回滚整个最后批次的所有迁移，而不是只回滚单个
- **批次追踪**：可以清楚地看到每个迁移是在哪个批次中执行的

### ✅ 事务保护

- 每个迁移在独立事务中执行
- 失败自动回滚，不会留下脏数据
- 确保迁移历史记录的一致性

### ✅ 版本控制

- 基于时间戳的版本号（YYYYMMDDHHMMSS）
- 自动按版本号排序
- 防止重复执行

## 快速开始

### 1. 使用 PowerShell 脚本（推荐）

```powershell
# 查看迁移状态
.\migrate.ps1 status

# 执行迁移
.\migrate.ps1 up

# 回滚最后一个批次
.\migrate.ps1 down
```

### 2. 使用 Makefile

```bash
# 查看迁移状态
make migrate-status

# 执行迁移
make migrate-up

# 回滚最后一个批次
make migrate-down
```

### 3. 直接使用 Go 命令

```bash
# 查看迁移状态
go run cmd/migrate/main.go -action=status

# 执行迁移
go run cmd/migrate/main.go -action=up

# 回滚最后一个批次
go run cmd/migrate/main.go -action=down

# 指定配置路径
go run cmd/migrate/main.go -action=up -config=config
```

## 迁移状态示例

```
迁移状态:
========================================
已执行的迁移:
  Batch 1:
    [✓] 20251103000000 - initial_schema
    [✓] 20251103000001 - add_indexes_for_performance
----------------------------------------
✓ 数据库已是最新，无待执行迁移
========================================
统计: 已执行 2 个, 待执行 0 个
```

## 创建新迁移

### Step 1: 生成版本号

```powershell
# Windows PowerShell
Get-Date -Format "yyyyMMddHHmmss"
# 输出示例: 20251103120000

# Linux/Mac
date +%Y%m%d%H%M%S
```

### Step 2: 创建迁移文件

在 `migrations/` 目录创建文件：

**文件名格式**: `YYYYMMDDHHMMSS_description.go`

例如：`migrations/20251103120000_add_user_avatar.go`

```go
package migrations

import (
	"gorm.io/gorm"
)

// AddUserAvatar20251103120000 添加用户头像字段
type AddUserAvatar20251103120000 struct{}

func (m *AddUserAvatar20251103120000) Version() string {
	return "20251103120000"
}

func (m *AddUserAvatar20251103120000) Name() string {
	return "add_user_avatar"
}

func (m *AddUserAvatar20251103120000) Up(db *gorm.DB) error {
	// 执行迁移
	return db.Exec(`
		ALTER TABLE users 
		ADD COLUMN avatar VARCHAR(255) COMMENT '头像URL' 
		AFTER nickname
	`).Error
}

func (m *AddUserAvatar20251103120000) Down(db *gorm.DB) error {
	// 回滚迁移
	return db.Exec("ALTER TABLE users DROP COLUMN avatar").Error
}
```

### Step 3: 注册迁移

编辑 `migrations/migrations.go`：

```go
func GetAllMigrations() []database.Migration {
	return []database.Migration{
		&Initial20251103000000{},
		&AddIndexes20251103000001{},
		&AddUserAvatar20251103120000{}, // 添加新迁移
	}
}
```

### Step 4: 执行迁移

```powershell
.\migrate.ps1 up
```

输出：
```
time="2025-11-03 12:00:00" level=info msg="执行迁移: 20251103120000 - add_user_avatar (batch: 2)"
time="2025-11-03 12:00:00" level=info msg="✓ 迁移完成: 20251103120000"
time="2025-11-03 12:00:00" level=info msg="所有迁移已完成 (共 1 个, batch: 2)"
```

## 迁移示例

### 添加字段

```go
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE users 
		ADD COLUMN phone VARCHAR(20) COMMENT '手机号' 
		AFTER email
	`).Error
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users DROP COLUMN phone").Error
}
```

### 创建表

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

### 添加索引

```go
func (m *Migration) Up(db *gorm.DB) error {
	// 检查索引是否存在
	var count int64
	db.Raw(`SELECT COUNT(*) FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users' 
		AND index_name = 'idx_email'`).Scan(&count)
	
	if count == 0 {
		return db.Exec("CREATE INDEX idx_email ON users(email)").Error
	}
	return nil
}

func (m *Migration) Down(db *gorm.DB) error {
	return db.Exec("DROP INDEX idx_email ON users").Error
}
```

### 数据迁移

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

## 批次管理示例

### 第一次部署（Batch 1）

```powershell
.\migrate.ps1 up
```

```
执行迁移: 20251103000000 - initial_schema (batch: 1)
✓ 迁移完成: 20251103000000
执行迁移: 20251103000001 - add_indexes_for_performance (batch: 1)
✓ 迁移完成: 20251103000001
所有迁移已完成 (共 2 个, batch: 1)
```

### 第二次部署（Batch 2）

添加3个新迁移后：

```powershell
.\migrate.ps1 up
```

```
执行迁移: 20251104000001 - add_user_avatar (batch: 2)
✓ 迁移完成: 20251104000001
执行迁移: 20251104000002 - add_notification_table (batch: 2)
✓ 迁移完成: 20251104000002
执行迁移: 20251104000003 - add_user_indexes (batch: 2)
✓ 迁移完成: 20251104000003
所有迁移已完成 (共 3 个, batch: 2)
```

### 回滚第二批（只回滚 Batch 2）

```powershell
.\migrate.ps1 down
```

```
开始回滚 batch 2 (共 3 个迁移)
回滚迁移: 20251104000003 - add_user_indexes
✓ 回滚完成: 20251104000003
回滚迁移: 20251104000002 - add_notification_table
✓ 回滚完成: 20251104000002
回滚迁移: 20251104000001 - add_user_avatar
✓ 回滚完成: 20251104000001
批次 2 回滚完成
```

**注意**：回滚是按照执行顺序的倒序进行的，确保依赖关系正确。

## 最佳实践

### 1. 迁移命名

✅ **Good**:
- `20251103000001_add_user_avatar.go`
- `20251103000002_create_notifications_table.go`
- `20251103000003_add_user_indexes.go`

❌ **Bad**:
- `migration1.go`
- `add_feature.go`
- `update.go`

### 2. 一次只做一件事

✅ **Good**: 一个迁移只添加一个字段或一个索引

❌ **Bad**: 一个迁移修改多个表和多个字段

### 3. 总是提供 Down() 方法

确保每个迁移都可以回滚：

```go
func (m *Migration) Down(db *gorm.DB) error {
	// ✅ Good: 具体的回滚逻辑
	return db.Exec("ALTER TABLE users DROP COLUMN avatar").Error
	
	// ❌ Bad: 空实现或 panic
	// return nil
	// panic("not implemented")
}
```

### 4. 在迁移中使用原生 SQL

避免依赖业务模型（模型可能会变）：

```go
// ✅ Good: 使用原生 SQL
func (m *Migration) Up(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users ADD COLUMN avatar VARCHAR(255)").Error
}

// ❌ Bad: 依赖业务模型
type User struct {
	// ... 字段可能会变化
}
func (m *Migration) Up(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
```

### 5. 测试迁移

在开发环境测试完整的 up/down 循环：

```powershell
# 1. 执行迁移
.\migrate.ps1 up

# 2. 验证数据库结构
# ... 检查表结构 ...

# 3. 回滚测试
.\migrate.ps1 down

# 4. 验证回滚
# ... 检查表结构是否恢复 ...

# 5. 再次执行
.\migrate.ps1 up
```

## 故障排除

### 问题 1: 迁移执行失败

**症状**: 迁移执行到一半失败

**原因**: 事务会自动回滚，不会记录到 migrations 表

**解决**:
1. 查看错误日志
2. 修复迁移代码
3. 重新执行 `migrate up`

### 问题 2: 迁移历史不一致

**症状**: 代码中的迁移与数据库记录不匹配

**解决**:
```sql
-- 查看当前记录
SELECT * FROM migrations ORDER BY batch, id;

-- 手动添加记录（如果手动执行了 SQL）
INSERT INTO migrations (version, migration, batch, created_at) 
VALUES ('20251103120000', 'add_user_avatar', 2, NOW());

-- 删除错误记录
DELETE FROM migrations WHERE version = '20251103120000';
```

### 问题 3: 批次号混乱

**原因**: 手动修改了 migrations 表

**解决**:
```sql
-- 重置批次号（谨慎使用！）
UPDATE migrations SET batch = 1 WHERE batch > 1;

-- 或者重新开始
TRUNCATE TABLE migrations;
-- 然后执行 migrate up
```

## 与现有系统集成

如果你有现有的 migrations 表：

1. **备份数据**
   ```sql
   CREATE TABLE migrations_backup AS SELECT * FROM migrations;
   ```

2. **检查表结构**
   ```sql
   DESC migrations;
   ```

3. **如果字段不匹配，调整表结构**
   ```sql
   ALTER TABLE migrations 
   ADD COLUMN batch INT(10) NOT NULL DEFAULT 0 COMMENT '批次号';
   
   ALTER TABLE migrations 
   CHANGE COLUMN migration migration VARCHAR(255) NOT NULL COMMENT '迁移名称';
   ```

4. **更新现有记录的批次号**
   ```sql
   UPDATE migrations SET batch = 1 WHERE batch = 0;
   ```

## 总结

✅ **实现的功能**:
- ✅ 基于时间戳的版本管理
- ✅ Batch 批次管理
- ✅ 事务保护
- ✅ Up/Down 双向迁移
- ✅ 迁移状态查看
- ✅ 与 Laravel/PHP migration 兼容的表结构

✅ **使用场景**:
- ✅ 团队协作开发
- ✅ 生产环境部署
- ✅ 数据库版本控制
- ✅ 一键回滚
- ✅ CI/CD 集成

🎉 **现在你可以像使用 PHP migration 一样管理 Go 项目的数据库了！**

