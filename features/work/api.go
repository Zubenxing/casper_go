package work

import (
	"casper_go/core/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ==================== Work APIs ====================

// CreateWorkAPI 创建工作记录
func CreateWorkAPI(c *gin.Context) {
	var req Work
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := CreateWork(userID, &req); err != nil {
		response.Error(c, 500, "创建失败: "+err.Error())
		return
	}

	response.Success(c, req)
}

// GetWorkListAPI 获取工作记录列表
func GetWorkListAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	works, total, err := GetWorkList(userID, status, page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  works,
		"total": total,
	})
}

// GetWorkDetailAPI 获取工作记录详情
func GetWorkDetailAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	work, err := GetWorkByID(userID, uint(id))
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, work)
}

// UpdateWorkAPI 更新工作记录
func UpdateWorkAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if err := UpdateWork(userID, uint(id), updates); err != nil {
		response.Error(c, 500, "更新失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteWorkAPI 删除工作记录
func DeleteWorkAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := DeleteWork(userID, uint(id)); err != nil {
		response.Error(c, 500, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetWorkStatsAPI 获取工作统计信息
func GetWorkStatsAPI(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := GetWorkStats(userID)
	if err != nil {
		response.Error(c, 500, "获取统计失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}

// ==================== WorkIssue APIs ====================

// CreateWorkIssueAPI 创建问题记录
func CreateWorkIssueAPI(c *gin.Context) {
	var req WorkIssue
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := CreateWorkIssue(userID, &req); err != nil {
		response.Error(c, 500, "创建失败: "+err.Error())
		return
	}

	response.Success(c, req)
}

// GetWorkIssueListAPI 获取问题记录列表
func GetWorkIssueListAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	issues, total, err := GetWorkIssueList(userID, status, page, pageSize)
	if err != nil {
		response.Error(c, 500, "获取失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  issues,
		"total": total,
	})
}

// GetWorkIssueDetailAPI 获取问题记录详情
func GetWorkIssueDetailAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	issue, err := GetWorkIssueByID(userID, uint(id))
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, issue)
}

// UpdateWorkIssueAPI 更新问题记录
func UpdateWorkIssueAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, "请求参数错误: "+err.Error())
		return
	}

	if err := UpdateWorkIssue(userID, uint(id), updates); err != nil {
		response.Error(c, 500, "更新失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteWorkIssueAPI 删除问题记录
func DeleteWorkIssueAPI(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := DeleteWorkIssue(userID, uint(id)); err != nil {
		response.Error(c, 500, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// GetWorkIssueStatsAPI 获取问题统计信息
func GetWorkIssueStatsAPI(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := GetWorkIssueStats(userID)
	if err != nil {
		response.Error(c, 500, "获取统计失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}
