package middleware

import (
	"strings"

	"casper_go/core/database"
	"casper_go/core/jwt"
	"casper_go/core/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			return
		}

		// 检查格式：Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证令牌格式错误")
			return
		}

		tokenString := parts[1]

		// 解析 token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "无效的认证令牌")
			return
		}

		// 验证 token 是否在数据库中且有效（避免循环依赖，直接查询）
		var count int64
		err = database.DB.Table("user_tokens").
			Where("token = ? AND is_valid = ? AND expires_at > NOW()", tokenString, true).
			Count(&count).Error
		if err != nil || count == 0 {
			response.Unauthorized(c, "认证令牌已失效")
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminAuth 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			return
		}
		c.Next()
	}
}

// GetUserID 获取当前用户 ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(uint)
	}
	return 0
}

// GetUsername 获取当前用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetUserRole 获取当前用户角色
func GetUserRole(c *gin.Context) string {
	if role, exists := c.Get("role"); exists {
		return role.(string)
	}
	return ""
}
