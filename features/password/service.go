package password

import (
	"casper_go/core/database"
	"casper_go/core/utils/crypto"
	"errors"

	"gorm.io/gorm"
)

// Create 创建账户
func Create(userID uint, req *CreateRequest) (*Account, error) {
	// 加密密码
	encryptedPassword, err := crypto.Encrypt(req.Password)
	if err != nil {
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
		return nil, err
	}

	return account, nil
}

// GetList 获取账户列表
func GetList(userID uint, page, pageSize int, category, keyword string) (*ListResponse, error) {
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
	account, err := GetByID(userID, id)
	if err != nil {
		return "", err
	}

	// 解密密码
	password, err := crypto.Decrypt(account.Password)
	if err != nil {
		return "", err
	}

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
	account, err := GetByID(userID, id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(account).Error; err != nil {
		return err
	}

	return nil
}

// GetCategories 获取用户的所有分类
func GetCategories(userID uint) ([]string, error) {
	var categories []string
	if err := database.DB.Model(&Account{}).
		Where("user_id = ? AND category != ''", userID).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetStats 获取统计信息
func GetStats(userID uint) (map[string]interface{}, error) {
	var total int64
	var favCount int64

	// 总数
	if err := database.DB.Model(&Account{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}

	// 收藏数
	if err := database.DB.Model(&Account{}).Where("user_id = ? AND is_fav = ?", userID, true).Count(&favCount).Error; err != nil {
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
		return nil, err
	}

	categoryMap := make(map[string]int64)
	for _, cc := range categoryCounts {
		categoryMap[cc.Category] = cc.Count
	}

	return map[string]interface{}{
		"total":           total,
		"fav_count":       favCount,
		"category_counts": categoryMap,
	}, nil
}
