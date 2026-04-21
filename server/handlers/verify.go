package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zihao-liu-qs/qs-blog/server/middleware"
	"github.com/zihao-liu-qs/qs-blog/server/services"
)

type verifyRequest struct {
	License  string `json:"license" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

// VerifyWithLogger 带日志的验证处理器
func VerifyWithLogger(logger *middleware.ActivityLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req verifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"valid": false, "reason": "bad_request"})
			return
		}

		result := services.Verify(req.License, req.DeviceID)

		// 异步记录日志（非阻塞）
		if logger != nil {
			reason := "success"
			if !result.Valid {
				reason = result.Reason
			}
			logger.Log(middleware.LogEntry{
				Time:     time.Now().Format(time.RFC3339),
				IP:       c.ClientIP(),
				License:  req.License,
				DeviceID: req.DeviceID,
				Success:  result.Valid,
				Reason:   reason,
			})
		}

		c.JSON(http.StatusOK, result)
	}
}
