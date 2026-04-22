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

	// Stripe 配置
	StripeKey        string // Stripe Secret Key
	StripePublishKey string // Stripe Publishable Key
	StripeWebhookSecret string // Stripe Webhook 签名密钥

	// 邮件配置
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string

	// 基础 URL
	BaseURL string
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
	logMaxAge := 30
	if v := os.Getenv("LOG_MAX_AGE_DAYS"); v != "" {
		if n, err := parseInt(v); err == nil {
			logMaxAge = n
		}
	}

	// Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	stripePublishKey := os.Getenv("STRIPE_PUBLISHABLE_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// Email
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if n, err := parseInt(v); err == nil {
			smtpPort = n
		}
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	fromName := os.Getenv("FROM_NAME")
	if fromName == "" {
		fromName = "QS Blog"
	}

	// Base URL
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &Config{
		Port: port, AdminKey: adminKey, DBPath: dbPath, LogPath: logPath, LogMaxAge: logMaxAge,
		StripeKey: stripeKey, StripePublishKey: stripePublishKey, StripeWebhookSecret: stripeWebhookSecret,
		SMTPHost: smtpHost, SMTPPort: smtpPort, SMTPUser: smtpUser, SMTPPassword: smtpPassword,
		FromEmail: fromEmail, FromName: fromName, BaseURL: baseURL,
	}
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
