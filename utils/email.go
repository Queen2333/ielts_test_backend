package utils

import (
	"fmt"
	"net/smtp"
)

func SendEmail(subject, body string, toEmail string) error {
	// 设置发件人的邮箱和密码
	from := "15502180712@163.com"
	password := "KDVEBGLNPYTTWZUC"

	// 设置 SMTP 服务器地址和端口
	smtpServer := "smtp.163.com"
	smtpPort := 25

	// 构建邮件内容
    msg := fmt.Sprintf("From: %s\r\n", from)
    msg += fmt.Sprintf("To: %s\r\n", toEmail)
    msg += fmt.Sprintf("Subject: %s\r\n", subject)
    msg += "\r\n" + body

    // SMTP认证信息
    auth := smtp.PlainAuth("", from, password, smtpServer)

	// 发送邮件
    err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, from, []string{toEmail}, []byte(msg))
    if err != nil {
        fmt.Println("邮件发送失败:", err)
        return err
    }
    fmt.Println("邮件发送成功！")

	return nil
}
