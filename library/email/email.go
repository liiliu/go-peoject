package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
	"your_project/library/logger"
)

// Config 邮件配置
type Config struct {
	Host     string // SMTP 服务器地址
	Port     int    // SMTP 端口
	Username string // 发件人邮箱
	Password string // 邮箱密码或授权码
	From     string // 发件人名称
}

var config Config

// Init 初始化邮件配置
func Init(cfg Config) {
	config = cfg
	logger.Info("email", "Email service initialized: %s:%d", cfg.Host, cfg.Port)
}

// Message 邮件消息
type Message struct {
	To          []string // 收件人列表
	Cc          []string // 抄送列表
	Bcc         []string // 密送列表
	Subject     string   // 主题
	Body        string   // 正文（支持 HTML）
	Attachments []string // 附件路径列表
}

// Send 发送邮件
func Send(msg Message) error {
	m := gomail.NewMessage()

	// 设置发件人
	m.SetHeader("From", fmt.Sprintf("%s <%s>", config.From, config.Username))

	// 设置收件人
	m.SetHeader("To", msg.To...)

	// 设置抄送
	if len(msg.Cc) > 0 {
		m.SetHeader("Cc", msg.Cc...)
	}

	// 设置密送
	if len(msg.Bcc) > 0 {
		m.SetHeader("Bcc", msg.Bcc...)
	}

	// 设置主题
	m.SetHeader("Subject", msg.Subject)

	// 设置正文（HTML）
	m.SetBody("text/html", msg.Body)

	// 添加附件
	for _, attachment := range msg.Attachments {
		m.Attach(attachment)
	}

	// 创建 SMTP 客户端
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		logger.Error("email", "Failed to send email: %v", err)
		return err
	}

	logger.Info("email", "Email sent successfully to %v", msg.To)
	return nil
}

// SendSimple 发送简单文本邮件
func SendSimple(to []string, subject, body string) error {
	return Send(Message{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

// SendHTML 发送 HTML 邮件
func SendHTML(to []string, subject, htmlBody string) error {
	return Send(Message{
		To:      to,
		Subject: subject,
		Body:    htmlBody,
	})
}

// ========================================
// 使用示例
// ========================================
/*
// 初始化
email.Init(email.Config{
    Host:     "smtp.gmail.com",
    Port:     587,
    Username: "your@gmail.com",
    Password: "your-password",
    From:     "Your Name",
})

// 发送简单邮件
email.SendSimple(
    []string{"recipient@example.com"},
    "Test Email",
    "This is a test email.",
)

// 发送 HTML 邮件
email.SendHTML(
    []string{"recipient@example.com"},
    "Welcome",
    "<h1>Welcome!</h1><p>Thanks for joining us.</p>",
)

// 发送带附件的邮件
email.Send(email.Message{
    To:          []string{"recipient@example.com"},
    Subject:     "Report",
    Body:        "<p>Please find the report attached.</p>",
    Attachments: []string{"/path/to/report.pdf"},
})
*/
