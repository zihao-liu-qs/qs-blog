package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zihao-liu-qs/qs-blog/server/services"
)

// CreateCheckout 创建 Stripe 支付会话
// POST /api/v1/checkout
func CreateCheckout(c *gin.Context) {
	var input services.CreateCheckoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文获取 Stripe Key（通过中间件注入）
	stripeKey := c.GetString("stripe_key")
	if stripeKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "stripe key not configured"})
		return
	}

	resp, err := services.CreateCheckoutSession(input, stripeKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
