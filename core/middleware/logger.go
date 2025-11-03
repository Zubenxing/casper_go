package middleware

import (
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

		logger.Log.WithFields(fields).Info("API 请求")
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
