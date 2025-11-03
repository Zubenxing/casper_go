package auth

import (
	"strings"

	"casper_go/core/middleware"
	"casper_go/core/response"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/login [post]
func LoginAPI(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	tokenResp, err := Login(req.Username, req.Password, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.SuccessMsg(c, "登录成功", tokenResp)
}

// Logout godoc
// @Summary 用户登出
// @Description 用户登出，使当前令牌失效
// @Tags 认证
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/logout [post]
func LogoutAPI(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		response.BadRequest(c, "认证令牌格式错误")
		return
	}

	if err := Logout(parts[1]); err != nil {
		response.ServerError(c, "登出失败: "+err.Error())
		return
	}

	response.SuccessMsg(c, "登出成功", nil)
}

// RefreshTokenAPI godoc
// @Summary 刷新令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body object{refresh_token=string} true "刷新令牌"
// @Success 200 {object} response.Response{data=TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/refresh [post]
func RefreshTokenAPI(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	tokenResp, err := RefreshToken(req.RefreshToken, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.SuccessMsg(c, "刷新成功", tokenResp)
}

// GetProfile godoc
// @Summary 获取当前用户信息
// @Description 获取已登录用户的详细信息
// @Tags 用户
// @Security BearerAuth
// @Success 200 {object} response.Response{data=UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/profile [get]
func GetProfileAPI(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := GetUserInfo(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}
