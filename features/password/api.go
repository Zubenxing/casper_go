package password

import (
	"casper_go/core/logger"
	"casper_go/core/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateAPI 创建账户
func CreateAPI(c *gin.Context) {
	// 从JWT中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Log.Infof("[密码管理] 用户 %v 创建账户: %s", userID, req.Title)

	account, err := Create(userID.(uint), &req)
	if err != nil {
		logger.Log.Errorf("[密码管理] 创建账户失败: %v", err)
		response.ServerError(c, "创建账户失败: "+err.Error())
		return
	}

	logger.Log.Infof("[密码管理] 账户创建成功: %s (ID: %d)", account.Title, account.ID)
	response.Success(c, account.ToResponse())
}

// GetListAPI 获取账户列表
func GetListAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	list, err := GetList(userID.(uint), page, pageSize, category, keyword)
	if err != nil {
		logger.Log.Errorf("[密码管理] 获取账户列表失败: %v", err)
		response.ServerError(c, "获取列表失败")
		return
	}

	response.Success(c, list)
}

// GetDetailAPI 获取账户详情（不含密码）
func GetDetailAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	account, err := GetByID(userID.(uint), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, account.ToResponse())
}

// GetPasswordAPI 获取密码（解密后）
func GetPasswordAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	logger.Log.Infof("[密码管理] 用户 %v 查看密码, 账户ID: %d", userID, id)

	password, err := GetPassword(userID.(uint), uint(id))
	if err != nil {
		response.ServerError(c, "获取密码失败: "+err.Error())
		return
	}

	response.Success(c, &PasswordResponse{
		Password: password,
	})
}

// UpdateAPI 更新账户
func UpdateAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	logger.Log.Infof("[密码管理] 用户 %v 更新账户, ID: %d", userID, id)

	account, err := Update(userID.(uint), uint(id), &req)
	if err != nil {
		logger.Log.Errorf("[密码管理] 更新账户失败: %v", err)
		response.ServerError(c, "更新账户失败: "+err.Error())
		return
	}

	logger.Log.Infof("[密码管理] 账户更新成功: %s (ID: %d)", account.Title, account.ID)
	response.Success(c, account.ToResponse())
}

// DeleteAPI 删除账户
func DeleteAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	logger.Log.Infof("[密码管理] 用户 %v 删除账户, ID: %d", userID, id)

	if err := Delete(userID.(uint), uint(id)); err != nil {
		logger.Log.Errorf("[密码管理] 删除账户失败: %v", err)
		response.ServerError(c, "删除账户失败: "+err.Error())
		return
	}

	logger.Log.Infof("[密码管理] 账户删除成功, ID: %d", id)
	response.Success(c, gin.H{"message": "删除成功"})
}

// GetCategoriesAPI 获取所有分类
func GetCategoriesAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	categories, err := GetCategories(userID.(uint))
	if err != nil {
		response.ServerError(c, "获取分类失败")
		return
	}

	response.Success(c, gin.H{
		"categories": categories,
	})
}

// GetStatsAPI 获取统计信息
func GetStatsAPI(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权")
		return
	}

	stats, err := GetStats(userID.(uint))
	if err != nil {
		response.ServerError(c, "获取统计信息失败")
		return
	}

	response.Success(c, stats)
}
