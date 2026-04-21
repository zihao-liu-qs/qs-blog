package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRateLimiter 按 IP 限流器
type IPRateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int           // 每分钟允许的请求次数
	window   time.Duration // 时间窗口
}

// Visitor 记录每个 IP 的访问信息
type Visitor struct {
	count      int
	firstVisit time.Time
}

// NewIPRateLimiter 创建限流器
func NewIPRateLimiter(rate int, window time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// 启动清理协程，每分钟清理过期记录
	go rl.cleanup()

	return rl
}

// Allow 检查是否允许请求
func (rl *IPRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &Visitor{
			count:      1,
			firstVisit: now,
		}
		return true
	}

	// 如果超过时间窗口，重置计数
	if now.Sub(v.firstVisit) > rl.window {
		v.count = 1
		v.firstVisit = now
		return true
	}

	// 在时间窗口内，检查是否超过限制
	v.count++
	return v.count <= rl.rate
}

// cleanup 定期清理过期的记录
func (rl *IPRateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, v := range rl.visitors {
			if now.Sub(v.firstVisit) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit 限流中间件
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded, try again later",
			})
			return
		}

		c.Next()
	}
}
