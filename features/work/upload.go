package work

import (
	"casper_go/config"
	"casper_go/core/logger"
	"casper_go/core/response"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadImage 上传问题截图
func UploadImage(c *gin.Context) {
	userID := c.GetUint("user_id")
	cfg := config.GlobalConfig.Upload

	logger.Work.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("[工作记录] 上传问题截图")

	file, err := c.FormFile("file")
	if err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 获取上传文件失败")
		response.BadRequest(c, "请选择要上传的图片")
		return
	}

	// 验证文件大小
	maxSize := cfg.GetMaxFileSizeBytes()
	if file.Size > maxSize {
		logger.Work.WithFields(logrus.Fields{
			"user_id":   userID,
			"file_size": file.Size,
		}).Warn("[工作记录] 上传文件过大")
		response.BadRequest(c, fmt.Sprintf("图片大小不能超过 %dMB", cfg.MaxFileSize))
		return
	}

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !strings.Contains(cfg.AllowedTypes, ext) {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"ext":     ext,
		}).Warn("[工作记录] 不支持的文件类型")
		response.BadRequest(c, "只支持图片格式: "+cfg.AllowedTypes)
		return
	}

	// 创建上传目录
	uploadDir := cfg.GetWorkIssuesPath()
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 创建上传目录失败")
		response.ServerError(c, "上传失败")
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := saveUploadedFile(file, filePath); err != nil {
		logger.Work.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("[工作记录] 保存文件失败")
		response.ServerError(c, "上传失败")
		return
	}

	// 只返回文件名，不返回完整URL（URL由前端拼接）
	logger.Work.WithFields(logrus.Fields{
		"user_id":  userID,
		"filename": filename,
	}).Info("[工作记录] 图片上传成功")

	response.Success(c, gin.H{
		"filename":     filename,      // 数据库存储的文件名
		"originalName": file.Filename, // 原始文件名
		"size":         file.Size,
	})
}

// saveUploadedFile 保存上传的文件
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
