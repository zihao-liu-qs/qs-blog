package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zihao-liu-qs/qs-blog/server/services"
)

func CreateLicense(c *gin.Context) {
	var input services.CreateLicenseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lic, err := services.CreateLicense(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, lic)
}

func ListLicenses(c *gin.Context) {
	list, err := services.ListLicenses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func GetLicense(c *gin.Context) {
	key := c.Param("key")
	detail, err := services.GetLicense(key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "license not found"})
		return
	}
	c.JSON(http.StatusOK, detail)
}

func RevokeLicense(c *gin.Context) {
	key := c.Param("key")
	if err := services.RevokeLicense(key); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "revoked"})
}

func RemoveDevice(c *gin.Context) {
	key := c.Param("key")
	deviceID := c.Param("device_id")
	if err := services.RemoveDeviceBinding(key, deviceID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}
