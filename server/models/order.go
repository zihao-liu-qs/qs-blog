package models

import "time"

// Order 订单模型
type Order struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	StripeSessionID string    `gorm:"uniqueIndex;not null"`
	StripePaymentID string    // Stripe 支付成功后的 ID
	ProductID       string    `gorm:"index;not null"` // 对应的软件产品 ID (如 blink)
	ProductName     string    // 产品名称 (如 Blink)
	ProductPrice    string    // 价格字符串 (如 ¥68)
	CustomerEmail   string    `gorm:"index;not null"`
	CustomerName    string
	LicenseKey      string // 支付成功后生成的 license key
	Status          string `gorm:"default:pending"` // pending, paid, failed
	Amount          int64  // Stripe 金额（分）
	Currency        string `gorm:"default:usd"`
	CreatedAt       time.Time
	PaidAt          *time.Time
}
