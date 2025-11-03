package middleware

import (
	"strings"
	"time"

	"casper_go/core/logger"
	"casper_go/core/response"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		fields := logrus.Fields{
			"status_code":  c.Writer.Status(),
			"latency_time": endTime.Sub(startTime),
			"client_ip":    c.ClientIP(),
			"req_method":   c.Request.Method,
			"req_uri":      c.Request.RequestURI,
		}

		// 添加请求ID（如果存在）
		if requestID := GetRequestID(c); requestID != "" {
			fields["request_id"] = requestID
		}

		// 根据请求路径，写入对应的模块日志
		uri := c.Request.RequestURI
		logged := false

		// 密码管理相关 API
		if strings.HasPrefix(uri, "/api/passwords") {
			logger.Password.WithFields(fields).Info("API 请求")
			logged = true
		}

		// 证书监控相关 API（包括证书列表、CSR、SSL验证等）
		if !logged && (strings.HasPrefix(uri, "/api/certificates") ||
			strings.HasPrefix(uri, "/api/csr/") ||
			strings.HasPrefix(uri, "/api/ssl-cert")) {
			logger.Certificate.WithFields(fields).Info("API 请求")
			logged = true
		}

		// 其他 API 写入通用日志
		if !logged {
			logger.Log.WithFields(fields).Info("API 请求")
		}
	}
}

// Recovery 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.WithFields(logrus.Fields{
					"error": err,
				}).Error("系统发生 panic")
				response.ServerError(c, "服务器内部错误")
			}
		}()
		c.Next()
	}
}
