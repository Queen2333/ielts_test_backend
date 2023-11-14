package utils

import "github.com/go-gomail/gomail"

func SendEmail(subject, body string, toEmail string) error {
	// 配置Gmail的SMTP服务器和端口
	d := gomail.NewDialer("smtp.gmail.com", 587, "breaknameyang@gmail.com", "xxef jkfc ndjx ajpm")

	// 创建邮件内容
	m := gomail.NewMessage()
	m.SetHeader("From", "breaknameyang@gmail.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
