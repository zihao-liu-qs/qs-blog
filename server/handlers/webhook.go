package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81/webhook"
	"github.com/zihao-liu-qs/qs-blog/server/services"
)

// StripeWebhook 处理 Stripe Webhook 回调
// POST /api/v1/webhook/stripe
func StripeWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	// 获取 Stripe 签名
	sigHeader := c.GetHeader("Stripe-Signature")
	endpointSecret := c.GetString("stripe_webhook_secret")

	// 验证签名
	event, err := webhook.ConstructEvent(body, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("webhook signature verification failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	// 处理不同事件
	switch event.Type {
	case "checkout.session.completed":
		var session map[string]interface{}
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("failed to parse session: %v", err)
			c.JSON(http.StatusOK, gin.H{"received": true})
			return
		}

		sessionID, _ := session["id"].(string)
		stripeKey := c.GetString("stripe_key")

		if err := services.HandlePaymentSuccess(sessionID, stripeKey); err != nil {
			log.Printf("failed to handle payment success: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("payment completed for session: %s", sessionID)

	default:
		log.Printf("unhandled event type: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
