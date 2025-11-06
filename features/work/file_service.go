package work

import (
	"casper_go/config"
	"casper_go/core/logger"
	"casper_go/core/response"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ServeFile 提供文件访问服务
// @Summary 获取上传的文件
// @Tags 文件服务
// @Param filename path string true "文件名（支持子目录，如 work-issues/xxx.png）"
// @Success 200 {file} binary "文件内容"
// @Failure 404 {object} response.Response "文件不存在"
// @Failure 403 {object} response.Response "禁止访问"
// @Router /api/files/{filename} [get]
func ServeFile(c *gin.Context) {
	filename := c.Param("filename")

	// Gin 的 *filename 参数会包含开头的斜杠，需要去除
	filename = strings.TrimPrefix(filename, "/")

	// 安全检查：防止路径穿越攻击
	if strings.Contains(filename, "..") {
		logger.Work.WithFields(logrus.Fields{
			"filename": filename,
			"ip":       c.ClientIP(),
		}).Warn("[文件服务] 检测到非法文件路径访问")
		response.Forbidden(c, "非法的文件路径")
		return
	}

	// 构建完整文件路径
	baseDir := config.GlobalConfig.Upload.BaseDir
	fullPath := filepath.Join(baseDir, filename)

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Work.WithFields(logrus.Fields{
			"filename": filename,
			"path":     fullPath,
		}).Warn("[文件服务] 文件不存在")
		response.NotFound(c, "文件不存在")
		return
	}

	// 记录文件访问日志
	logger.Work.WithFields(logrus.Fields{
		"filename": filename,
		"fullPath": fullPath,
		"ip":       c.ClientIP(),
	}).Info("[文件服务] 文件访问成功")

	// 设置合适的 Content-Type
	ext := strings.ToLower(filepath.Ext(filename))
	contentType := getContentType(ext)
	if contentType != "" {
		c.Header("Content-Type", contentType)
	}

	// 设置缓存头（图片可以缓存1小时）
	c.Header("Cache-Control", "public, max-age=3600")

	// 返回文件
	c.File(fullPath)
}

// getContentType 根据文件扩展名返回 Content-Type
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".bmp":  "image/bmp",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}

// GetFileURL 生成文件访问URL
func GetFileURL(filename string) string {
	// 移除开头的斜杠（如果有）
	filename = strings.TrimPrefix(filename, "/")
	filename = strings.TrimPrefix(filename, "uploads/")

	return fmt.Sprintf("/api/files/%s", filename)
}
