package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"github.com/zihao-liu-qs/qs-blog/server/database"
	"github.com/zihao-liu-qs/qs-blog/server/models"
)

type VerifyResult struct {
	Valid       bool       `json:"valid"`
	Reason      string     `json:"reason,omitempty"`
	ExpireAt    *time.Time `json:"expire_at,omitempty"`
	DeviceCount int        `json:"device_count,omitempty"`
	MaxDevices  int        `json:"max_devices,omitempty"`
}

func Verify(licenseKey, deviceID string) VerifyResult {
	var lic models.License
	if err := database.DB.Where("key = ?", licenseKey).First(&lic).Error; err != nil {
		return VerifyResult{Valid: false, Reason: "invalid_license"}
	}
	if !lic.IsActive {
		return VerifyResult{Valid: false, Reason: "revoked"}
	}
	if lic.ExpireAt != nil && time.Now().After(*lic.ExpireAt) {
		return VerifyResult{Valid: false, Reason: "expired"}
	}

	var binding models.DeviceBinding
	err := database.DB.Where("license_key = ? AND device_id = ?", licenseKey, deviceID).First(&binding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var count int64
		database.DB.Model(&models.DeviceBinding{}).Where("license_key = ?", licenseKey).Count(&count)
		if int(count) >= lic.MaxDevices {
			return VerifyResult{Valid: false, Reason: "device_limit_reached"}
		}
		now := time.Now()
		binding = models.DeviceBinding{LicenseKey: licenseKey, DeviceID: deviceID, FirstSeenAt: now, LastSeenAt: now}
		database.DB.Create(&binding)
	} else {
		database.DB.Model(&binding).Update("last_seen_at", time.Now())
	}

	var deviceCount int64
	database.DB.Model(&models.DeviceBinding{}).Where("license_key = ?", licenseKey).Count(&deviceCount)
	return VerifyResult{Valid: true, ExpireAt: lic.ExpireAt, DeviceCount: int(deviceCount), MaxDevices: lic.MaxDevices}
}

type CreateLicenseInput struct {
	Key        string     `json:"key" binding:"required"`
	AppID      string     `json:"app_id"`
	MaxDevices int        `json:"max_devices"`
	ExpireAt   *time.Time `json:"expire_at"`
	Note       string     `json:"note"`
}

func CreateLicense(input CreateLicenseInput) (*models.License, error) {
	if input.MaxDevices <= 0 {
		input.MaxDevices = 1
	}
	lic := models.License{
		Key: input.Key, AppID: input.AppID, MaxDevices: input.MaxDevices,
		ExpireAt: input.ExpireAt, Note: input.Note, IsActive: true,
	}
	if err := database.DB.Create(&lic).Error; err != nil {
		return nil, err
	}
	return &lic, nil
}

func RevokeLicense(key string) error {
	result := database.DB.Model(&models.License{}).Where("key = ?", key).Update("is_active", false)
	if result.RowsAffected == 0 {
		return errors.New("license not found")
	}
	return result.Error
}

type LicenseDetail struct {
	models.License
	BoundDevices []models.DeviceBinding `json:"bound_devices"`
}

func GetLicense(key string) (*LicenseDetail, error) {
	var lic models.License
	if err := database.DB.Where("key = ?", key).First(&lic).Error; err != nil {
		return nil, err
	}
	var bindings []models.DeviceBinding
	database.DB.Where("license_key = ?", key).Find(&bindings)
	return &LicenseDetail{License: lic, BoundDevices: bindings}, nil
}

func ListLicenses() ([]models.License, error) {
	var list []models.License
	err := database.DB.Order("id desc").Find(&list).Error
	return list, err
}

func RemoveDeviceBinding(licenseKey, deviceID string) error {
	result := database.DB.Where("license_key = ? AND device_id = ?", licenseKey, deviceID).Delete(&models.DeviceBinding{})
	if result.RowsAffected == 0 {
		return errors.New("binding not found")
	}
	return result.Error
}
