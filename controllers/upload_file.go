package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 上传文件
// @Description 上传一个文件到服务器
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} map[string]string{"message": "string"}
// @Failure 400 {object} map[string]string{"message": "string"}
// @Failure 500 {object} map[string]string{"message": "string"}
// @Router /upload [post]
func UploadFile(c *gin.Context) {
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 定义上传文件的路径
	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "upload file err: %s")
		return
	}

	baseURL := fmt.Sprintf("%s://%s", c.Request.URL.Scheme, c.Request.Host)
	fileURL := strings.Join([]string{baseURL, "uploads", filename}, "/")

	utils.HandleResponse(c, http.StatusOK, fileURL, "Success")
}