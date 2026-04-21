package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	AdminKey  string
	DBPath    string
	LogPath   string
	LogMaxAge int // 日志保留天数
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		adminKey = "change-me-in-production"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/activate.db"
	}
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		logPath = "./data/activity.log"
	}
	logMaxAge := 30 // 默认保留 30 天
	if v := os.Getenv("LOG_MAX_AGE_DAYS"); v != "" {
		if n, err := parseInt(v); err == nil {
			logMaxAge = n
		}
	}
	return &Config{Port: port, AdminKey: adminKey, DBPath: dbPath, LogPath: logPath, LogMaxAge: logMaxAge}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
