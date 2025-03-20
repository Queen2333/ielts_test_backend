package controllers

import (
	"net/http"

	"github.com/Queen2333/ielts_test_backend/services"
	"github.com/gin-gonic/gin"
)

// 其他已有代码保持不变...

// AskGrok 处理用户向 Grok 提问的请求
func AskGrok(c *gin.Context) {
	var request struct {
		Question string `json:"question" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "问题不能为空"})
		return
	}

	// 调用 Grok API
	response, err := services.CallGrokAPI(request.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "调用Grok API失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": response})
}