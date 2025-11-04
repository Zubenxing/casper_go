package auth

import (
	"errors"
	"time"

	"casper_go/config"
	"casper_go/core/database"
	"casper_go/core/jwt"

	"gorm.io/gorm"
)

// Login 用户登录
func Login(username, password, ipAddress, userAgent string) (*TokenResponse, error) {
	var user User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	if !jwt.CheckPassword(password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	token, err := jwt.GenerateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	cfg := config.GlobalConfig.JWT
	now := time.Now()

	userToken := &UserToken{
		UserID:           user.UserID,
		Token:            token,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        now.Add(time.Duration(cfg.ExpireHours) * time.Hour),
		RefreshExpiresAt: now.Add(time.Duration(cfg.RefreshExpireHours) * time.Hour),
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		IsValid:          true,
	}

	if err := database.DB.Create(userToken).Error; err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    userToken.ExpiresAt,
		User:         user.ToResponse(),
	}, nil
}

// Logout 用户登出
func Logout(token string) error {
	return database.DB.Model(&UserToken{}).
		Where("token = ?", token).
		Update("is_valid", false).Error
}

// RefreshToken 刷新令牌
func RefreshToken(refreshToken, ipAddress, userAgent string) (*TokenResponse, error) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	var userToken UserToken
	if err := database.DB.Where("refresh_token = ? AND is_valid = ?", refreshToken, true).First(&userToken).Error; err != nil {
		return nil, errors.New("刷新令牌已失效")
	}

	if userToken.IsRefreshExpired() {
		return nil, errors.New("刷新令牌已过期")
	}

	var user User
	if err := database.DB.Where("user_id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	newToken, err := jwt.GenerateToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(user.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	cfg := config.GlobalConfig.JWT
	now := time.Now()

	database.DB.Model(&userToken).Update("is_valid", false)

	newUserToken := &UserToken{
		UserID:           user.UserID,
		Token:            newToken,
		RefreshToken:     newRefreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        now.Add(time.Duration(cfg.ExpireHours) * time.Hour),
		RefreshExpiresAt: now.Add(time.Duration(cfg.RefreshExpireHours) * time.Hour),
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		IsValid:          true,
	}

	if err := database.DB.Create(newUserToken).Error; err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  newToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresAt:    newUserToken.ExpiresAt,
		User:         user.ToResponse(),
	}, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(userID uint) (*UserResponse, error) {
	var user User
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

// ValidateToken 验证 Token 是否有效
func ValidateToken(token string) bool {
	var userToken UserToken
	err := database.DB.Where("token = ? AND is_valid = ?", token, true).First(&userToken).Error
	if err != nil {
		return false
	}
	return !userToken.IsExpired()
}

// InitDefaultData 初始化默认数据
func InitDefaultData() error {
	var count int64
	database.DB.Model(&User{}).Where("username = ?", "admin").Count(&count)

	if count == 0 {
		hashedPassword, err := jwt.HashPassword("admin123")
		if err != nil {
			return err
		}

		// 获取当前最大的 admin user_id，如果没有则从 1000 开始
		var maxUserID uint = 999
		database.DB.Model(&User{}).Where("role = ? AND user_id >= ? AND user_id < ?", "admin", 1000, 2000).
			Select("COALESCE(MAX(user_id), 999)").Scan(&maxUserID)

		admin := &User{
			UserID:   maxUserID + 1,
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@casper.local",
			Nickname: "系统管理员",
			Role:     "admin",
			Status:   1,
		}

		if err := database.DB.Create(admin).Error; err != nil {
			return err
		}
	}

	return nil
}
