package models

import "time"

type License struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Key        string    `gorm:"uniqueIndex;not null"`
	AppID      string    `gorm:"index"`
	MaxDevices int       `gorm:"default:1"`
	ExpireAt   *time.Time
	Note       string
	IsActive   bool      `gorm:"default:true"`
	CreatedAt  time.Time
}

type DeviceBinding struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	LicenseKey  string    `gorm:"index;not null"`
	DeviceID    string    `gorm:"not null"`
	FirstSeenAt time.Time
	LastSeenAt  time.Time
}
