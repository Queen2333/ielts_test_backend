package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 发送验证码到邮箱
// @Description 生成6位随机验证码并发送到指定邮箱地址
// @Accept json
// @Produce json
// @Param email body string true "目标邮箱地址"
// @Success 200 {object} models.ResponseData
// @Failure 400 {object} models.ResponseData
// @Failure 400 {object} models.ResponseData
// @Failure 500 {object} models.ResponseData
// @Router /send-code [post]
func SendCodeHandler(c *gin.Context) {
	utils.InitRedis()

	// 解析请求参数
	var request struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
        return
    }

	// 验证邮箱格式
	if err := utils.IsValidEmail(request.Email); !err {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid email")
        return
	}

	// 生成6位随机验证码
	code := utils.GenerateRandomNumber(6)

	// 将验证码存储到 Redis 中，有效期为5分钟
	err := utils.Set("verification_code_" + request.Email, code, 15 * time.Minute)
    if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to store verification code")
        return
    }

	utils.SendEmail("邮箱登录验证", fmt.Sprintf("您的验证码为: %s, 有效期为15分钟", code), request.Email)

	// 返回验证码
	utils.HandleResponse(c, http.StatusOK, code, "success")
}
