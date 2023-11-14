package utils

import "github.com/gin-gonic/gin"

// GetClientIP 获取用户IP
func GetClientIP(c *gin.Context) string {
	if c.Request.Header.Get("X-Forwarded-For") != "" {
		return c.Request.Header.Get("X-Forwarded-For")
	}
	return c.ClientIP()
}
