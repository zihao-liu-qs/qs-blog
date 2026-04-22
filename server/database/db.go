package database

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/zihao-liu-qs/qs-blog/server/models"
)

var DB *gorm.DB

func Init(dbPath string) error {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return err
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.License{}, &models.DeviceBinding{}, &models.Order{}); err != nil {
		return err
	}
	DB = db
	return nil
}
