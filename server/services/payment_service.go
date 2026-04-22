package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/zihao-liu-qs/qs-blog/server/database"
	"github.com/zihao-liu-qs/qs-blog/server/models"
)

// CreateCheckoutInput 创建支付会话的输入
type CreateCheckoutInput struct {
	ProductID     string `json:"product_id" binding:"required"`
	ProductName   string `json:"product_name" binding:"required"`
	ProductPrice  int64  `json:"product_price" binding:"required"` // 价格（分），如 ¥68 = 6800 分（CNY）
	CustomerEmail string `json:"customer_email" binding:"required"`
	CustomerName  string `json:"customer_name"`
	SuccessURL    string `json:"success_url"`
	CancelURL     string `json:"cancel_url"`
	Currency      string `json:"currency"` // usd, cny 等
}

// CheckoutResponse 创建支付会话的响应
type CheckoutResponse struct {
	SessionID      string `json:"session_id"`
	PublishableKey string `json:"publishable_key"`
	URL            string `json:"url"` // Stripe 跳转链接
}

// CreateCheckoutSession 创建 Stripe Checkout Session
func CreateCheckoutSession(input CreateCheckoutInput, stripeKey string) (*CheckoutResponse, error) {
	stripe.Key = stripeKey

	if input.Currency == "" {
		input.Currency = "cny"
	}

	// 创建订单项
	lineItems := []*stripe.CheckoutSessionLineItemParams{
		{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(input.Currency),
				UnitAmount: stripe.Int64(input.ProductPrice),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(input.ProductName),
				},
			},
			Quantity: stripe.Int64(1),
		},
	}

	// 创建 Stripe Checkout Session
	params := &stripe.CheckoutSessionParams{
		LineItems: lineItems,
		Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(input.SuccessURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(input.CancelURL),
		CustomerEmail: stripe.String(input.CustomerEmail),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"product_id":   input.ProductID,
				"product_name": input.ProductName,
			},
		},
		Metadata: map[string]string{
			"product_id":    input.ProductID,
			"product_name":  input.ProductName,
			"customer_name": input.CustomerName,
		},
	}

	s, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("stripe session creation failed: %w", err)
	}

	// 保存订单到数据库
	order := models.Order{
		StripeSessionID: s.ID,
		ProductID:       input.ProductID,
		ProductName:     input.ProductName,
		ProductPrice:    formatPrice(input.ProductPrice, input.Currency),
		CustomerEmail:   input.CustomerEmail,
		CustomerName:    input.CustomerName,
		Status:          "pending",
		Amount:          input.ProductPrice,
		Currency:        input.Currency,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		log.Printf("failed to save order: %v", err)
	}

	return &CheckoutResponse{
		SessionID:      s.ID,
		PublishableKey: stripeKey,
		URL:            s.URL,
	}, nil
}

// HandlePaymentSuccess 支付成功后的处理（Webhook 调用）
func HandlePaymentSuccess(sessionID, stripeKey string) error {
	stripe.Key = stripeKey

	// 从 Stripe 获取会话详情
	s, err := session.Get(sessionID, &stripe.CheckoutSessionParams{
		Expand: []*string{stripe.String("payment_intent"), stripe.String("line_items")},
	})
	if err != nil {
		return fmt.Errorf("failed to get stripe session: %w", err)
	}

	if s.PaymentStatus != stripe.CheckoutSessionPaymentStatusPaid {
		return fmt.Errorf("payment not completed")
	}

	// 查找本地订单
	var order models.Order
	if err := database.DB.Where("stripe_session_id = ?", sessionID).First(&order).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	if order.Status == "paid" {
		// 已处理过，幂等
		return nil
	}

	// 更新订单状态
	now := time.Now()
	order.Status = "paid"
	if s.PaymentIntent != nil {
		order.StripePaymentID = s.PaymentIntent.ID
	}
	order.PaidAt = &now

	// 生成 License Key
	licenseKey := GenerateLicenseKey()
	order.LicenseKey = licenseKey

	if err := database.DB.Save(&order).Error; err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	// 创建 License 记录
	_, err = CreateLicense(CreateLicenseInput{
		Key:        licenseKey,
		AppID:      order.ProductID,
		MaxDevices: 1,
		Note:       fmt.Sprintf("Auto-generated for order %s, customer: %s", order.StripeSessionID, order.CustomerEmail),
	})
	if err != nil {
		return fmt.Errorf("failed to create license: %w", err)
	}

	// 发送邮件
	go SendLicenseEmail(order.CustomerEmail, order.CustomerName, order.ProductName, licenseKey, order.ProductPrice)

	return nil
}

// GenerateLicenseKey 生成唯一的 License Key
func GenerateLicenseKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	part1 := hex.EncodeToString(b)[:12]

	b2 := make([]byte, 4)
	rand.Read(b2)
	part2 := hex.EncodeToString(b2)[:8]

	return fmt.Sprintf("%s-%s", part1, part2)
}

func formatPrice(price int64, currency string) string {
	if currency == "cny" {
		return fmt.Sprintf("¥%.2f", float64(price)/100)
	}
	if currency == "usd" {
		return fmt.Sprintf("$%.2f", float64(price)/100)
	}
	return fmt.Sprintf("%d", price)
}
