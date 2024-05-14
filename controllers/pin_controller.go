package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// 发送验证码接口
func SendCodeHandler(c *gin.Context) {
	utils.InitRedis()

	// 解析请求参数
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

	// 验证邮箱格式
	if err := utils.IsValidEmail(request.Email); !err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
        return
	}

	// 生成6位随机验证码
	code := utils.GenerateRandomNumber(6)

	// 将验证码存储到 Redis 中，有效期为5分钟
	err := utils.Set("verification_code_" + request.Email, code, 15 * time.Minute)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store verification code"})
        return
    }

	utils.SendEmail("邮箱登录验证", fmt.Sprintf("您的验证码为: %s, 有效期为15分钟", code), request.Email)

	// 返回验证码
	c.JSON(http.StatusOK, gin.H{"code": code})
}
