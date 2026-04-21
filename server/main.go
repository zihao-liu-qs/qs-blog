package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zihao-liu-qs/qs-blog/server/config"
	"github.com/zihao-liu-qs/qs-blog/server/database"
	"github.com/zihao-liu-qs/qs-blog/server/handlers"
	"github.com/zihao-liu-qs/qs-blog/server/middleware"
)

func main() {
	cfg := config.Load()

	if err := database.Init(cfg.DBPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// 确保日志目录存在
	os.MkdirAll("./data", 0755)

	// 创建活动日志记录器
	activityLogger, err := middleware.NewActivityLogger(cfg.LogPath, time.Duration(cfg.LogMaxAge)*24*time.Hour)
	if err != nil {
		log.Fatalf("failed to init activity logger: %v", err)
	}
	defer activityLogger.Close()

	r := gin.Default()

	// 创建限流器：每个 IP 每分钟最多 5 次请求
	rateLimiter := middleware.NewIPRateLimiter(5, time.Minute)

	// 公开接口（启用限流 + 日志）
	api := r.Group("/api/v1", middleware.RateLimit(rateLimiter))
	{
		api.POST("/activate", handlers.VerifyWithLogger(activityLogger))
	}

	// 管理接口（需要 X-Admin-Key）
	admin := r.Group("/api/admin", middleware.AdminAuth(cfg.AdminKey))
	{
		admin.POST("/licenses", handlers.CreateLicense)
		admin.GET("/licenses", handlers.ListLicenses)
		admin.GET("/licenses/:key", handlers.GetLicense)
		admin.DELETE("/licenses/:key", handlers.RevokeLicense)
		admin.DELETE("/devices/:key/:device_id", handlers.RemoveDevice)
	}

	log.Printf("server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
