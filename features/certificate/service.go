package certificate

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"casper_go/config"
	"casper_go/core/database"
	"casper_go/core/logger"

	"gorm.io/gorm"
)

// Add 添加证书监控
func Add(url string) (*Monitor, error) {
	// 先检查是否已存在（包括已软删除的）
	var existing Monitor
	err := database.DB.Unscoped().Where("url = ?", url).First(&existing).Error

	if err == nil {
		// URL 已存在
		if existing.DeletedAt.Valid {
			// 如果是软删除的，恢复它
			logger.Log.Infof("[证书服务] 恢复已删除的 URL: %s", url)
			existing.DeletedAt = gorm.DeletedAt{}
			database.DB.Unscoped().Save(&existing)
			return &existing, nil
		}
		logger.Log.Warnf("[证书服务] URL 已存在: %s", url)
		return nil, errors.New("该 URL 已存在监控列表")
	}

	// URL 不存在，开始检查证书
	logger.Log.Infof("[证书服务] 开始检查证书: %s", url)
	certInfo, err := Check(url)
	if err != nil {
		logger.Log.Errorf("[证书服务] 证书检查失败 URL=%s, 错误: %v", url, err)
		return nil, fmt.Errorf("证书检查失败: %w", err)
	}

	cfg := config.GlobalConfig.Certificate
	status := GetStatus(certInfo.DaysLeft, cfg.WarningDays)

	monitor := &Monitor{
		URL:          url,
		Domain:       certInfo.Domain,
		Issuer:       certInfo.Issuer,
		Subject:      certInfo.Subject,
		Organization: certInfo.Organization,
		NotBefore:    certInfo.NotBefore,
		NotAfter:     certInfo.NotAfter,
		DaysLeft:     certInfo.DaysLeft,
		IsValid:      certInfo.IsValid,
		ErrorMsg:     certInfo.Error,
		LastCheckAt:  time.Now(),
		Status:       status,
	}

	if err := database.DB.Create(monitor).Error; err != nil {
		logger.Log.Errorf("[证书服务] 数据库插入失败 URL=%s, 错误: %v", url, err)
		return nil, fmt.Errorf("数据库错误: %w", err)
	}

	logger.Log.Infof("[证书服务] 证书添加成功 URL=%s, ID=%d, 剩余%d天", url, monitor.ID, monitor.DaysLeft)
	return monitor, nil
}

// GetList 获取监控列表
func GetList(page, pageSize int) ([]Monitor, int64, error) {
	var monitors []Monitor
	var total int64

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&Monitor{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := database.DB.Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&monitors).Error; err != nil {
		return nil, 0, err
	}

	return monitors, total, nil
}

// GetByID 根据 ID 获取监控信息
func GetByID(id uint) (*Monitor, error) {
	var monitor Monitor
	if err := database.DB.First(&monitor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("监控记录不存在")
		}
		return nil, err
	}
	return &monitor, nil
}

// Update 更新证书监控信息
func Update(id uint) (*Monitor, error) {
	monitor, err := GetByID(id)
	if err != nil {
		return nil, err
	}

	certInfo, err := Check(monitor.URL)
	if err != nil {
		monitor.ErrorMsg = err.Error()
		monitor.IsValid = false
		monitor.LastCheckAt = time.Now()
		database.DB.Save(monitor)
		return nil, err
	}

	cfg := config.GlobalConfig.Certificate
	status := GetStatus(certInfo.DaysLeft, cfg.WarningDays)

	monitor.Domain = certInfo.Domain
	monitor.Issuer = certInfo.Issuer
	monitor.Subject = certInfo.Subject
	monitor.Organization = certInfo.Organization
	monitor.NotBefore = certInfo.NotBefore
	monitor.NotAfter = certInfo.NotAfter
	monitor.DaysLeft = certInfo.DaysLeft
	monitor.IsValid = certInfo.IsValid
	monitor.ErrorMsg = certInfo.Error
	monitor.LastCheckAt = time.Now()
	monitor.Status = status

	if err := database.DB.Save(monitor).Error; err != nil {
		return nil, err
	}

	return monitor, nil
}

// Delete 删除证书监控
func Delete(id uint) error {
	return database.DB.Delete(&Monitor{}, id).Error
}

// CheckAllResult 批量检查结果
type CheckAllResult struct {
	Total     int       `json:"total"`      // 总数
	Success   int       `json:"success"`    // 成功数
	Failed    int       `json:"failed"`     // 失败数
	Duration  string    `json:"duration"`   // 耗时
	StartTime time.Time `json:"start_time"` // 开始时间
	EndTime   time.Time `json:"end_time"`   // 结束时间
}

// CheckAll 并发检查所有监控的证书
func CheckAll() (*CheckAllResult, error) {
	return CheckAllWithContext(context.Background())
}

// CheckAllWithContext 带超时的并发检查所有证书
func CheckAllWithContext(ctx context.Context) (*CheckAllResult, error) {
	startTime := time.Now()
	
	// 获取所有需要检查的证书（排除删除的）
	var monitors []Monitor
	if err := database.DB.Find(&monitors).Error; err != nil {
		return nil, fmt.Errorf("查询证书列表失败: %w", err)
	}

	total := len(monitors)
	if total == 0 {
		logger.Log.Info("[证书服务] 没有需要检查的证书")
		return &CheckAllResult{
			Total:     0,
			Success:   0,
			Failed:    0,
			StartTime: startTime,
			EndTime:   time.Now(),
			Duration:  "0s",
		}, nil
	}

	logger.Log.Infof("[证书服务] 开始并发检查 %d 个证书", total)

	// 使用并发处理，限制并发数为 10
	concurrency := 10
	if total < concurrency {
		concurrency = total
	}

	// 创建任务通道和结果通道
	taskChan := make(chan Monitor, total)
	resultChan := make(chan bool, total)

	// 启动工作协程池
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for monitor := range taskChan {
				select {
				case <-ctx.Done():
					logger.Log.Warnf("[证书服务] Worker %d 收到取消信号", workerID)
					resultChan <- false
					return
				default:
					// 执行证书检查
					logger.Log.Debugf("[证书服务] Worker %d 检查: %s", workerID, monitor.URL)
					_, err := updateMonitor(&monitor)
					resultChan <- (err == nil)
				}
			}
		}(i)
	}

	// 发送任务到通道
	go func() {
		for _, monitor := range monitors {
			taskChan <- monitor
		}
		close(taskChan)
	}()

	// 等待所有工作协程完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 统计结果
	success := 0
	failed := 0
	for result := range resultChan {
		if result {
			success++
		} else {
			failed++
		}
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	result := &CheckAllResult{
		Total:     total,
		Success:   success,
		Failed:    failed,
		StartTime: startTime,
		EndTime:   endTime,
		Duration:  duration.Round(time.Millisecond).String(),
	}

	logger.Log.Infof("[证书服务] 检查完成: 总数=%d, 成功=%d, 失败=%d, 耗时=%s", 
		result.Total, result.Success, result.Failed, result.Duration)

	return result, nil
}

// updateMonitor 更新单个证书监控（内部方法）
func updateMonitor(monitor *Monitor) (*Monitor, error) {
	certInfo, err := Check(monitor.URL)
	
	// 无论成功失败都更新最后检查时间
	monitor.LastCheckAt = time.Now()
	
	if err != nil {
		monitor.ErrorMsg = err.Error()
		monitor.IsValid = false
		// 保存错误信息
		if saveErr := database.DB.Save(monitor).Error; saveErr != nil {
			logger.Log.Errorf("[证书服务] 保存错误信息失败 URL=%s: %v", monitor.URL, saveErr)
		}
		return nil, err
	}

	cfg := config.GlobalConfig.Certificate
	status := GetStatus(certInfo.DaysLeft, cfg.WarningDays)

	// 更新证书信息
	monitor.Domain = certInfo.Domain
	monitor.Issuer = certInfo.Issuer
	monitor.Subject = certInfo.Subject
	monitor.Organization = certInfo.Organization
	monitor.NotBefore = certInfo.NotBefore
	monitor.NotAfter = certInfo.NotAfter
	monitor.DaysLeft = certInfo.DaysLeft
	monitor.IsValid = certInfo.IsValid
	monitor.ErrorMsg = certInfo.Error
	monitor.Status = status

	if err := database.DB.Save(monitor).Error; err != nil {
		logger.Log.Errorf("[证书服务] 保存证书信息失败 URL=%s: %v", monitor.URL, err)
		return nil, err
	}

	return monitor, nil
}
