package middleware

import (
	"sync"
	"time"

	"casper_go/core/errors"
	"casper_go/core/response"

	"github.com/gin-gonic/gin"
)

// 基于内存的简单限流器
type rateLimiter struct {
	requests map[string]*requestInfo
	mu       sync.RWMutex
	limit    int           // 每分钟最大请求数
	window   time.Duration // 时间窗口
}

type requestInfo struct {
	count     int
	resetTime time.Time
}

var limiter *rateLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter(requestsPerMinute int) {
	limiter = &rateLimiter{
		requests: make(map[string]*requestInfo),
		limit:    requestsPerMinute,
		window:   time.Minute,
	}

	// 定期清理过期数据
	go limiter.cleanup()
}

// cleanup 清理过期的限流记录
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, info := range rl.requests {
			if now.After(info.resetTime) {
				delete(rl.requests, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		// 使用 IP 作为限流 key
		key := c.ClientIP()

		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		now := time.Now()
		info, exists := limiter.requests[key]

		if !exists || now.After(info.resetTime) {
			// 新的时间窗口
			limiter.requests[key] = &requestInfo{
				count:     1,
				resetTime: now.Add(limiter.window),
			}
			c.Next()
			return
		}

		if info.count >= limiter.limit {
			// 超过限流阈值
			c.Header("X-RateLimit-Limit", string(rune(limiter.limit)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", string(rune(info.resetTime.Unix())))

			response.Error(c, errors.ErrTooManyRequest, errors.GetMessage(errors.ErrTooManyRequest))
			c.Abort()
			return
		}

		// 增加计数
		info.count++
		c.Header("X-RateLimit-Limit", string(rune(limiter.limit)))
		c.Header("X-RateLimit-Remaining", string(rune(limiter.limit-info.count)))

		c.Next()
	}
}
