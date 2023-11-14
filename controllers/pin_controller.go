package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lux208716/go-gin-project/utils"
	"net/http"
	"os"
)

// SendPinController 处理发送PIN码的请求
func SendPinController(c *gin.Context) {
	// 解析请求中的JSON数据
	var requestBody struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !utils.IsValidEmail(requestBody.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入正确的邮箱地址"})
		return
	}

	/*************************** 注册发送邮件 ********************************/
	// 模板字符串包含一个PIN码的占位符
	templateString := `<!DOCTYPE html> 
    <html>
      <head>
      <title>PIN码页面</title>
	  </head>
	  <body>
		<h1>Welcome to the PIN Code Page</h1>
		<p style="margin: 0 10px;">Your PIN code is: {{.PINCode}}</p>
        <p>验证码15分钟内有效！</p>
	  </body>
	</html>
	`
	// 生成HTML页面并传入PIN码和模板
	pinCode := utils.GeneratePIN(6) // 生成一个6位的随机PIN码
	htmlPage, err := utils.GenerateHTMLPageWithPIN(pinCode, templateString)
	if err != nil {
		fmt.Println("Error generating HTML page:", err)
		os.Exit(1)
	}
	utils.SendEmail("[Go server]", htmlPage, requestBody.Email)

	// 模拟发送PIN码给用户邮箱（实际中需要使用真实的邮件服务）
	// 这里只是简单地打印出来
	fmt.Printf("Sending PIN code %s to email: %s\n", pinCode, requestBody.Email)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "PIN code sent successfully"})
}
