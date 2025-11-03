package certificate

import (
	"casper_go/core/database"
	"casper_go/core/logger"
)

// UpdateInfo 更新证书的客户名和备注
func UpdateInfo(id uint, customerName, remark string) (*Monitor, error) {
	monitor, err := GetByID(id)
	if err != nil {
		return nil, err
	}

	monitor.CustomerName = customerName
	monitor.Remark = remark

	if err := database.DB.Save(monitor).Error; err != nil {
		logger.Log.Errorf("[证书服务] 更新信息失败 ID=%d, 错误: %v", id, err)
		return nil, err
	}

	logger.Log.Infof("[证书服务] 更新信息成功 ID=%d, 客户名=%s", id, customerName)
	return monitor, nil
}
