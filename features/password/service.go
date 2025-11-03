package password

import (
	"casper_go/core/database"
	"casper_go/core/logger"
	"casper_go/core/utils/crypto"
	"errors"

	"gorm.io/gorm"
)

// Create 创建账户
func Create(userID uint, req *CreateRequest) (*Account, error) {
	logger.Password.Infof("[密码管理] 开始创建账户: 用户ID=%d, 标题=%s", userID, req.Title)

	// 加密密码
	encryptedPassword, err := crypto.Encrypt(req.Password)
	if err != nil {
		logger.Password.Errorf("[密码管理] 密码加密失败: 用户ID=%d, 错误=%v", userID, err)
		return nil, err
	}

	account := &Account{
		UserID:   userID,
		Title:    req.Title,
		URL:      req.URL,
		Username: req.Username,
		Password: encryptedPassword,
		Category: req.Category,
		Notes:    req.Notes,
		IsFav:    req.IsFav,
	}

	if err := database.DB.Create(account).Error; err != nil {
		logger.Password.Errorf("[密码管理] 创建账户失败: 用户ID=%d, 错误=%v", userID, err)
		return nil, err
	}

	logger.Password.Infof("[密码管理] 创建账户成功: ID=%d, 标题=%s", account.ID, account.Title)
	return account, nil
}

// GetList 获取账户列表
func GetList(userID uint, page, pageSize int, category, keyword string) (*ListResponse, error) {
	logger.Password.Infof("[密码管理] 获取账户列表: 用户ID=%d, 页码=%d, 每页=%d, 分类=%s, 关键词=%s",
		userID, page, pageSize, category, keyword)

	var accounts []Account
	var total int64

	query := database.DB.Model(&Account{}).Where("user_id = ?", userID)

	// 按分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("title LIKE ? OR url LIKE ? OR username LIKE ? OR notes LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询（按创建时间正序，收藏的在前）
	offset := (page - 1) * pageSize
	if err := query.Order("is_fav DESC, created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&accounts).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]*Response, len(accounts))
	for i, account := range accounts {
		items[i] = account.ToResponse()
	}

	logger.Password.Infof("[密码管理] 获取账户列表成功: 总数=%d, 返回=%d条", total, len(items))

	return &ListResponse{
		Total: total,
		Items: items,
	}, nil
}

// GetByID 根据ID获取账户
func GetByID(userID, id uint) (*Account, error) {
	var account Account
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账户不存在")
		}
		return nil, err
	}
	return &account, nil
}

// GetPassword 获取解密后的密码
func GetPassword(userID, id uint) (string, error) {
	logger.Password.Infof("[密码管理] 获取密码: 用户ID=%d, 账户ID=%d", userID, id)

	account, err := GetByID(userID, id)
	if err != nil {
		logger.Password.Warnf("[密码管理] 账户不存在: 用户ID=%d, 账户ID=%d", userID, id)
		return "", err
	}

	// 解密密码
	password, err := crypto.Decrypt(account.Password)
	if err != nil {
		logger.Password.Errorf("[密码管理] 密码解密失败: 账户ID=%d, 错误=%v", id, err)
		return "", err
	}

	logger.Password.Infof("[密码管理] 获取密码成功: 账户ID=%d, 标题=%s", id, account.Title)
	return password, nil
}

// Update 更新账户
func Update(userID, id uint, req *UpdateRequest) (*Account, error) {
	account, err := GetByID(userID, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})

	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		// 加密新密码
		encryptedPassword, err := crypto.Encrypt(req.Password)
		if err != nil {
			return nil, err
		}
		updates["password"] = encryptedPassword
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.IsFav != nil {
		updates["is_fav"] = *req.IsFav
	}

	if err := database.DB.Model(account).Updates(updates).Error; err != nil {
		return nil, err
	}

	return account, nil
}

// Delete 删除账户
func Delete(userID, id uint) error {
	logger.Password.Infof("[密码管理] 开始删除账户: 用户ID=%d, 账户ID=%d", userID, id)

	account, err := GetByID(userID, id)
	if err != nil {
		logger.Password.Warnf("[密码管理] 账户不存在: 用户ID=%d, 账户ID=%d", userID, id)
		return err
	}

	if err := database.DB.Delete(account).Error; err != nil {
		logger.Password.Errorf("[密码管理] 删除账户失败: 账户ID=%d, 错误=%v", id, err)
		return err
	}

	logger.Password.Infof("[密码管理] 删除账户成功: ID=%d, 标题=%s", id, account.Title)
	return nil
}

// GetCategories 获取用户的所有分类
func GetCategories(userID uint) ([]string, error) {
	logger.Password.Infof("[密码管理] 获取分类列表: 用户ID=%d", userID)

	var categories []string
	if err := database.DB.Model(&Account{}).
		Where("user_id = ? AND category != ''", userID).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		logger.Password.Errorf("[密码管理] 获取分类列表失败: %v", err)
		return nil, err
	}

	logger.Password.Infof("[密码管理] 获取分类列表成功: 共%d个分类", len(categories))
	return categories, nil
}

// GetStats 获取统计信息（优化版：减少查询次数）
func GetStats(userID uint) (map[string]interface{}, error) {
	logger.Password.Infof("[密码管理] 获取统计信息: 用户ID=%d", userID)

	// 使用一次查询获取总数和收藏数
	type CountResult struct {
		Total    int64
		FavCount int64
	}

	var countResult CountResult
	if err := database.DB.Model(&Account{}).
		Select("COUNT(*) as total, SUM(CASE WHEN is_fav = 1 THEN 1 ELSE 0 END) as fav_count").
		Where("user_id = ?", userID).
		Scan(&countResult).Error; err != nil {
		logger.Password.Errorf("[密码管理] 获取统计信息失败: %v", err)
		return nil, err
	}

	// 按分类统计
	type CategoryCount struct {
		Category string
		Count    int64
	}
	var categoryCounts []CategoryCount
	if err := database.DB.Model(&Account{}).
		Select("category, COUNT(*) as count").
		Where("user_id = ? AND category != ''", userID).
		Group("category").
		Scan(&categoryCounts).Error; err != nil {
		logger.Password.Errorf("[密码管理] 获取分类统计失败: %v", err)
		return nil, err
	}

	categoryMap := make(map[string]int64)
	for _, cc := range categoryCounts {
		categoryMap[cc.Category] = cc.Count
	}

	logger.Password.Infof("[密码管理] 统计信息获取成功: 总数=%d, 收藏=%d, 分类数=%d",
		countResult.Total, countResult.FavCount, len(categoryMap))

	return map[string]interface{}{
		"total":           countResult.Total,
		"fav_count":       countResult.FavCount,
		"category_counts": categoryMap,
	}, nil
}
