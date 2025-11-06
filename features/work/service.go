package work

import (
	"casper_go/core/database"
	"casper_go/core/logger"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ==================== Work Service ====================

// CreateWork 创建工作记录
func CreateWork(userID uint, work *Work) error {
	logger.Work.WithFields(logrus.Fields{
		"user_id":  userID,
		"title":    work.Title,
		"status":   work.Status,
		"priority": work.Priority,
	}).Info("[工作记录] 创建工作记录")

	work.UserID = userID
	if err := database.DB.Create(work).Error; err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 创建失败")
		return err
	}

	logger.Work.WithFields(logrus.Fields{
		"work_id":  work.ID,
		"title":    work.Title,
		"status":   work.Status,
		"priority": work.Priority,
	}).Info("[工作记录] 创建成功")
	return nil
}

// GetWorkList 获取工作记录列表
func GetWorkList(userID uint, status string, page, pageSize int) ([]Work, int64, error) {
	logger.Work.WithFields(logrus.Fields{
		"user_id":   userID,
		"status":    status,
		"page":      page,
		"page_size": pageSize,
	}).Info("[工作记录] 获取列表")

	var works []Work
	var total int64

	query := database.DB.Model(&Work{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 获取总数失败")
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&works).Error; err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 获取列表失败")
		return nil, 0, err
	}

	logger.Work.WithFields(logrus.Fields{
		"user_id": userID,
		"total":   total,
		"count":   len(works),
	}).Info("[工作记录] 获取列表成功")
	return works, total, nil
}

// GetWorkByID 根据ID获取工作记录
func GetWorkByID(userID, id uint) (*Work, error) {
	logger.Work.Infof("[工作管理] 获取工作详情: ID=%d, 用户ID=%d", id, userID)

	var work Work
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&work).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Work.Warnf("[工作管理] 工作记录不存在: ID=%d", id)
			return nil, errors.New("工作记录不存在")
		}
		logger.Work.Errorf("[工作管理] 获取工作详情失败: %v", err)
		return nil, err
	}

	logger.Work.Infof("[工作管理] 获取工作详情成功: ID=%d", id)
	return &work, nil
}

// UpdateWork 更新工作记录
func UpdateWork(userID, id uint, updates map[string]interface{}) error {
	logger.Work.Infof("[工作管理] 更新工作记录: ID=%d, 用户ID=%d", id, userID)

	// 检查记录是否存在
	work, err := GetWorkByID(userID, id)
	if err != nil {
		return err
	}

	// 处理日期字段：将ISO 8601格式转换为time.Time
	dateFields := []string{"start_date", "due_date"}
	for _, field := range dateFields {
		if dateStr, ok := updates[field].(string); ok && dateStr != "" {
			// 解析 ISO 8601 格式: 2025-11-05T16:00:00.000Z
			parsedTime, err := time.Parse(time.RFC3339, dateStr)
			if err != nil {
				logger.Work.Warnf("[工作管理] 日期解析失败 %s=%s: %v", field, dateStr, err)
				// 如果解析失败，尝试其他格式
				parsedTime, err = time.Parse("2006-01-02T15:04:05Z", dateStr)
				if err != nil {
					logger.Work.Errorf("[工作管理] 日期格式错误 %s=%s", field, dateStr)
					return errors.New("日期格式错误: " + field)
				}
			}
			updates[field] = parsedTime
		}
	}

	// 如果状态变更为 completed，自动设置完成时间
	if status, ok := updates["status"].(string); ok && status == WorkStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := database.DB.Model(work).Updates(updates).Error; err != nil {
		logger.Work.Errorf("[工作管理] 更新工作记录失败: %v", err)
		return err
	}

	logger.Work.Infof("[工作管理] 更新工作记录成功: ID=%d", id)
	return nil
}

// DeleteWork 删除工作记录
func DeleteWork(userID, id uint) error {
	logger.Work.Infof("[工作管理] 删除工作记录: ID=%d, 用户ID=%d", id, userID)

	work, err := GetWorkByID(userID, id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(work).Error; err != nil {
		logger.Work.Errorf("[工作管理] 删除工作记录失败: %v", err)
		return err
	}

	logger.Work.Infof("[工作管理] 删除工作记录成功: ID=%d", id)
	return nil
}

// GetWorkStats 获取工作统计信息
func GetWorkStats(userID uint) (map[string]interface{}, error) {
	logger.Work.Infof("[工作管理] 获取工作统计: 用户ID=%d", userID)

	var total, pending, inProgress, completed int64

	// 一次查询获取所有统计
	type StatusCount struct {
		Total      int64
		Pending    int64
		InProgress int64
		Completed  int64
	}

	var result StatusCount
	err := database.DB.Model(&Work{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as completed
		`, WorkStatusPending, WorkStatusInProgress, WorkStatusCompleted).
		Where("user_id = ?", userID).
		Scan(&result).Error

	if err != nil {
		logger.Work.Errorf("[工作管理] 获取工作统计失败: %v", err)
		return nil, err
	}

	total = result.Total
	pending = result.Pending
	inProgress = result.InProgress
	completed = result.Completed

	stats := map[string]interface{}{
		"total":       total,
		"pending":     pending,
		"in_progress": inProgress,
		"completed":   completed,
	}

	logger.Work.Infof("[工作管理] 统计信息获取成功: 总数=%d, 待办=%d, 进行中=%d, 已完成=%d",
		total, pending, inProgress, completed)
	return stats, nil
}

// ==================== WorkIssue Service ====================

// CreateWorkIssue 创建问题记录
func CreateWorkIssue(userID uint, issue *WorkIssue) error {
	logger.Work.Infof("[问题管理] 创建问题记录: 用户ID=%d, 标题=%s", userID, issue.Title)

	issue.UserID = userID

	// 处理 images 字段：空字符串转为空（MySQL JSON 不接受空字符串）
	if issue.Images == "" {
		issue.Images = "[]" // 设置为空 JSON 数组
	}

	if err := database.DB.Create(issue).Error; err != nil {
		logger.Work.Errorf("[问题管理] 创建问题记录失败: %v", err)
		return err
	}

	logger.Work.Infof("[问题管理] 创建问题记录成功: ID=%d", issue.ID)
	return nil
}

// GetWorkIssueList 获取问题记录列表
func GetWorkIssueList(userID uint, status string, page, pageSize int) ([]WorkIssue, int64, error) {
	logger.Work.Infof("[问题管理] 获取问题列表: 用户ID=%d, 状态=%s, 页码=%d, 每页=%d", userID, status, page, pageSize)

	var issues []WorkIssue
	var total int64

	query := database.DB.Model(&WorkIssue{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Work.Errorf("[问题管理] 获取问题总数失败: %v", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&issues).Error; err != nil {
		logger.Work.Errorf("[问题管理] 获取问题列表失败: %v", err)
		return nil, 0, err
	}

	logger.Work.Infof("[问题管理] 获取问题列表成功: 总数=%d, 返回=%d条", total, len(issues))
	return issues, total, nil
}

// GetWorkIssueByID 根据ID获取问题记录
func GetWorkIssueByID(userID, id uint) (*WorkIssue, error) {
	logger.Work.Infof("[问题管理] 获取问题详情: ID=%d, 用户ID=%d", id, userID)

	var issue WorkIssue
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&issue).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Work.Warnf("[问题管理] 问题记录不存在: ID=%d", id)
			return nil, errors.New("问题记录不存在")
		}
		logger.Work.Errorf("[问题管理] 获取问题详情失败: %v", err)
		return nil, err
	}

	logger.Work.Infof("[问题管理] 获取问题详情成功: ID=%d", id)
	return &issue, nil
}

// UpdateWorkIssue 更新问题记录
func UpdateWorkIssue(userID, id uint, updates map[string]interface{}) error {
	logger.Work.Infof("[问题管理] 更新问题记录: ID=%d, 用户ID=%d", id, userID)

	// 检查记录是否存在
	issue, err := GetWorkIssueByID(userID, id)
	if err != nil {
		return err
	}

	// 处理 images 字段：空字符串转为空数组（MySQL JSON 不接受空字符串）
	if images, ok := updates["images"].(string); ok {
		if images == "" {
			updates["images"] = "[]" // 设置为空 JSON 数组
		}
	}

	// 如果状态变更为 resolved，自动设置解决时间
	if status, ok := updates["status"].(string); ok && status == IssueStatusResolved {
		now := time.Now()
		updates["resolved_at"] = &now
	}

	if err := database.DB.Model(issue).Updates(updates).Error; err != nil {
		logger.Work.Errorf("[问题管理] 更新问题记录失败: %v", err)
		return err
	}

	logger.Work.Infof("[问题管理] 更新问题记录成功: ID=%d", id)
	return nil
}

// DeleteWorkIssue 删除问题记录
func DeleteWorkIssue(userID, id uint) error {
	logger.Work.Infof("[问题管理] 删除问题记录: ID=%d, 用户ID=%d", id, userID)

	issue, err := GetWorkIssueByID(userID, id)
	if err != nil {
		return err
	}

	if err := database.DB.Delete(issue).Error; err != nil {
		logger.Work.Errorf("[问题管理] 删除问题记录失败: %v", err)
		return err
	}

	logger.Work.Infof("[问题管理] 删除问题记录成功: ID=%d", id)
	return nil
}

// GetWorkIssueStats 获取问题统计信息
func GetWorkIssueStats(userID uint) (map[string]interface{}, error) {
	logger.Work.Infof("[问题管理] 获取问题统计: 用户ID=%d", userID)

	var total, open, inProgress, resolved int64

	// 一次查询获取所有统计
	type StatusCount struct {
		Total      int64
		Open       int64
		InProgress int64
		Resolved   int64
	}

	var result StatusCount
	err := database.DB.Model(&WorkIssue{}).
		Select(`
			COUNT(*) as total,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as open,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as in_progress,
			SUM(CASE WHEN status = ? THEN 1 ELSE 0 END) as resolved
		`, IssueStatusOpen, IssueStatusInProgress, IssueStatusResolved).
		Where("user_id = ?", userID).
		Scan(&result).Error

	if err != nil {
		logger.Work.Errorf("[问题管理] 获取问题统计失败: %v", err)
		return nil, err
	}

	total = result.Total
	open = result.Open
	inProgress = result.InProgress
	resolved = result.Resolved

	stats := map[string]interface{}{
		"total":       total,
		"open":        open,
		"in_progress": inProgress,
		"resolved":    resolved,
	}

	logger.Work.Infof("[问题管理] 统计信息获取成功: 总数=%d, 待处理=%d, 进行中=%d, 已解决=%d",
		total, open, inProgress, resolved)
	return stats, nil
}
