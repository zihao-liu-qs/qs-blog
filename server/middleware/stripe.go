package middleware

import "github.com/gin-gonic/gin"

// InjectStripeKey 将 Stripe Key 注入到上下文
func InjectStripeKey(stripeKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("stripe_key", stripeKey)
		c.Next()
	}
}

// InjectStripeWebhookSecret 将 Webhook Secret 注入到上下文
func InjectStripeWebhookSecret(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("stripe_webhook_secret", secret)
		c.Next()
	}
}
