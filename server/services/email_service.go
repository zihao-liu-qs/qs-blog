package services

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// EmailConfig 邮件服务配置
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

var emailCfg *EmailConfig

// InitEmail 初始化邮件配置
func InitEmail(cfg *EmailConfig) {
	emailCfg = cfg
}

// SendLicenseEmail 发送 License 邮件
func SendLicenseEmail(customerEmail, customerName, productName, licenseKey, price string) {
	if emailCfg == nil {
		log.Println("email config not initialized, skipping")
		return
	}

	subject := fmt.Sprintf("感谢您的购买！%s 的 License Key", productName)

	htmlBody := fmt.Sprintf(`
<html>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
  <h2 style="color: #333;">您好，%s！</h2>
  <p>感谢您购买 <strong>%s</strong>。</p>
  
  <div style="background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0;">
    <p style="margin: 0; color: #666;">您的 License Key：</p>
    <p style="font-size: 24px; font-family: monospace; margin: 10px 0; color: #333; letter-spacing: 2px;">%s</p>
  </div>

  <h3>激活步骤：</h3>
  <ol style="color: #555;">
    <li>打开 %s 软件</li>
    <li>在激活界面输入上方的 License Key</li>
    <li>点击"激活"按钮</li>
  </ol>

  <p style="color: #888; font-size: 14px;">
    订单金额：%s<br>
    如有问题，请回复此邮件联系我们。
  </p>
  
  <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
  <p style="color: #aaa; font-size: 12px;">此邮件由系统自动发送，请勿直接回复。</p>
</body>
</html>
`, customerName, productName, licenseKey, productName, price)

	textBody := fmt.Sprintf(`
您好，%s！

感谢您购买 %s。

您的 License Key：%s

激活步骤：
1. 打开 %s 软件
2. 在激活界面输入上方的 License Key
3. 点击"激活"按钮

订单金额：%s
如有问题，请回复此邮件联系我们。
`, customerName, productName, licenseKey, productName, price)

	e := &email.Email{
		To:      []string{customerEmail},
		From:    fmt.Sprintf("%s <%s>", emailCfg.FromName, emailCfg.FromEmail),
		Subject: subject,
		HTML:    []byte(htmlBody),
		Text:    []byte(textBody),
	}

	addr := fmt.Sprintf("%s:%d", emailCfg.SMTPHost, emailCfg.SMTPPort)
	auth := smtp.PlainAuth("", emailCfg.SMTPUser, emailCfg.SMTPPassword, emailCfg.SMTPHost)

	if err := e.Send(addr, auth); err != nil {
		log.Printf("failed to send email to %s: %v", customerEmail, err)
	} else {
		log.Printf("license email sent to %s successfully", customerEmail)
	}
}
