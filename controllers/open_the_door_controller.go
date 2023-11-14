package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os/exec"
)

// ExecutePythonScriptHandler 是执行Python脚本的控制器
func ExecutePythonScriptHandler(c *gin.Context) {
	cmd := exec.Command("python3", "scripts/python/open_the_door.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": string(output)})
}
