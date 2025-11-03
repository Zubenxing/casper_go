package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

type AddUserIDField20251103000005 struct{}

func (m *AddUserIDField20251103000005) Version() string {
	return "20251103000005"
}

func (m *AddUserIDField20251103000005) Name() string {
	return "add_user_id_field"
}

func (m *AddUserIDField20251103000005) Up(tx *gorm.DB) error {
	// 1. 检查 user_id 字段是否已存在
	var hasUserID bool
	if err := tx.Raw(`
		SELECT COUNT(*) > 0 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users' 
		AND column_name = 'user_id'
	`).Scan(&hasUserID).Error; err != nil {
		return err
	}

	// 如果字段不存在，则添加
	if !hasUserID {
		if err := tx.Exec(`
			ALTER TABLE users 
			ADD COLUMN user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID(1000-2000:admin, 2000-10000:普通用户)' 
			AFTER id
		`).Error; err != nil {
			return fmt.Errorf("添加 user_id 字段失败: %w", err)
		}
	}

	// 2. 为 user_id 为 0 的用户生成 user_id（避免重复执行）
	// admin 用户从 1000 开始
	if err := tx.Exec(`
		SET @row_number = (SELECT COALESCE(MAX(user_id), 999) FROM users WHERE user_id >= 1000 AND user_id < 2000) - 999
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		UPDATE users 
		SET user_id = 1000 + (@row_number := @row_number + 1)
		WHERE role = 'admin' AND user_id = 0
		ORDER BY id
	`).Error; err != nil {
		return fmt.Errorf("为 admin 用户生成 user_id 失败: %w", err)
	}

	// 普通用户从 2000 开始
	if err := tx.Exec(`
		SET @row_number = (SELECT COALESCE(MAX(user_id), 1999) FROM users WHERE user_id >= 2000) - 1999
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		UPDATE users 
		SET user_id = 2000 + (@row_number := @row_number + 1)
		WHERE role != 'admin' AND user_id = 0
		ORDER BY id
	`).Error; err != nil {
		return fmt.Errorf("为普通用户生成 user_id 失败: %w", err)
	}

	// 3. 检查并添加唯一索引
	var hasIndex bool
	if err := tx.Raw(`
		SELECT COUNT(*) > 0
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users' 
		AND index_name = 'idx_users_user_id'
	`).Scan(&hasIndex).Error; err != nil {
		return err
	}

	if !hasIndex {
		if err := tx.Exec(`
			CREATE UNIQUE INDEX idx_users_user_id ON users(user_id)
		`).Error; err != nil {
			return fmt.Errorf("创建 user_id 唯一索引失败: %w", err)
		}
	}

	// 4. 删除外键约束（如果存在）
	// 检查外键是否存在
	var hasForeignKey bool
	tx.Raw(`
		SELECT COUNT(*) > 0
		FROM information_schema.table_constraints
		WHERE constraint_schema = DATABASE()
		AND table_name = 'user_tokens'
		AND constraint_name = 'fk_users_tokens'
		AND constraint_type = 'FOREIGN KEY'
	`).Scan(&hasForeignKey)

	if hasForeignKey {
		tx.Exec("ALTER TABLE user_tokens DROP FOREIGN KEY fk_users_tokens")
	}

	// 5. 更新其他表的 user_id 字段以匹配新的 user_id
	// 只更新那些 user_id 对应 users.id 的记录

	// 更新 user_tokens 表
	if err := tx.Exec(`
		UPDATE user_tokens ut
		INNER JOIN users u ON ut.user_id = u.id
		SET ut.user_id = u.user_id
		WHERE u.user_id != 0
	`).Error; err != nil {
		return fmt.Errorf("更新 user_tokens.user_id 失败: %w", err)
	}

	// 更新 password_accounts 表
	if err := tx.Exec(`
		UPDATE password_accounts pa
		INNER JOIN users u ON pa.user_id = u.id
		SET pa.user_id = u.user_id
		WHERE u.user_id != 0
	`).Error; err != nil {
		return fmt.Errorf("更新 password_accounts.user_id 失败: %w", err)
	}

	// 更新 works 表
	if err := tx.Exec(`
		UPDATE works w
		INNER JOIN users u ON w.user_id = u.id
		SET w.user_id = u.user_id
		WHERE u.user_id != 0
	`).Error; err != nil {
		return fmt.Errorf("更新 works.user_id 失败: %w", err)
	}

	// 更新 work_issues 表
	if err := tx.Exec(`
		UPDATE work_issues wi
		INNER JOIN users u ON wi.user_id = u.id
		SET wi.user_id = u.user_id
		WHERE u.user_id != 0
	`).Error; err != nil {
		return fmt.Errorf("更新 work_issues.user_id 失败: %w", err)
	}

	return nil
}

func (m *AddUserIDField20251103000005) Down(tx *gorm.DB) error {
	// 1. 恢复其他表的 user_id 为 users.id
	if err := tx.Exec(`
		UPDATE user_tokens ut
		INNER JOIN users u ON ut.user_id = u.user_id
		SET ut.user_id = u.id
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		UPDATE password_accounts pa
		INNER JOIN users u ON pa.user_id = u.user_id
		SET pa.user_id = u.id
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		UPDATE works w
		INNER JOIN users u ON w.user_id = u.user_id
		SET w.user_id = u.id
	`).Error; err != nil {
		return err
	}

	if err := tx.Exec(`
		UPDATE work_issues wi
		INNER JOIN users u ON wi.user_id = u.user_id
		SET wi.user_id = u.id
	`).Error; err != nil {
		return err
	}

	// 2. 删除索引（检查是否存在）
	var hasIndex bool
	tx.Raw(`
		SELECT COUNT(*) > 0
		FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users' 
		AND index_name = 'idx_users_user_id'
	`).Scan(&hasIndex)

	if hasIndex {
		tx.Exec("DROP INDEX idx_users_user_id ON users")
	}

	// 3. 删除 user_id 字段（检查是否存在）
	var hasUserID bool
	tx.Raw(`
		SELECT COUNT(*) > 0 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = 'users' 
		AND column_name = 'user_id'
	`).Scan(&hasUserID)

	if hasUserID {
		if err := tx.Exec("ALTER TABLE users DROP COLUMN user_id").Error; err != nil {
			return err
		}
	}

	return nil
}
