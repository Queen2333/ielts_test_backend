package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 允许的文件类型
var allowedImageTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

var allowedAudioTypes = map[string]bool{
	".mp3": true,
	".wav": true,
	".m4a": true,
	".ogg": true,
}

// @Summary 上传文件（Go原生）
// @Description 上传图片或音频文件，直接保存到本地
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Param type query string false "文件类型：image 或 audio"
// @Success 200 {object} models.ResponseData{data=map[string]string}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /upload/file [post]
func UploadFileNative(c *gin.Context) {
	// 获取上传的文件（支持多种字段名）
	file, err := c.FormFile("file")
	if err != nil {
		// 尝试使用 "audio" 字段名
		file, err = c.FormFile("audio")
		if err != nil {
			// 尝试使用 "image" 字段名
			file, err = c.FormFile("image")
			if err != nil {
				utils.HandleResponse(c, http.StatusBadRequest, nil, "No file uploaded")
				return
			}
		}
	}

	// 检查文件大小（限制为 10MB）
	maxSize := int64(10 * 1024 * 1024) // 10MB
	if file.Size > maxSize {
		utils.HandleResponse(c, http.StatusBadRequest, nil, "File size exceeds 10MB limit")
		return
	}

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 自动识别文件类型（根据扩展名）
	var fileType string
	if allowedImageTypes[ext] {
		fileType = "image"
	} else if allowedAudioTypes[ext] {
		fileType = "audio"
	} else {
		// 如果无法识别，使用查询参数
		fileType = c.DefaultQuery("type", "")
	}

	// 调试日志
	fmt.Printf("[Upload Debug] Filename: %s, Extension: %s, Auto-detected FileType: %s\n", file.Filename, ext, fileType)

	// 验证文件类型并设置上传目录
	var uploadDir string

	if fileType == "image" {
		uploadDir = "uploads/images"
	} else if fileType == "audio" {
		uploadDir = "uploads/audio"
	} else {
		utils.HandleResponse(c, http.StatusBadRequest, nil, fmt.Sprintf("Unsupported file type: %s", ext))
		return
	}

	// 创建上传目录（如果不存在）
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, nil, "Failed to create upload directory")
		return
	}

	// 生成唯一文件名：日期_UUID_原文件名
	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()[:8]
	newFilename := fmt.Sprintf("%s_%s_%s", timestamp, uniqueID, file.Filename)
	filePath := filepath.Join(uploadDir, newFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, nil, "Failed to save file")
		return
	}

	// 返回文件访问 URL
	fileURL := fmt.Sprintf("/files/%s/%s", fileType+"s", newFilename)

	response := map[string]interface{}{
		"url":      fileURL,
		"filename": newFilename,
		"size":     file.Size,
		"type":     fileType,
	}

	utils.HandleResponse(c, http.StatusOK, response, "File uploaded successfully")
}

// @Summary 删除文件
// @Description 根据文件 URL 删除文件
// @Tags File
// @Accept json
// @Produce json
// @Param url body map[string]string true "文件 URL"
// @Success 200 {object} models.ResponseData{data=nil}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /upload/delete [delete]
func DeleteFile(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, nil, "Invalid request")
		return
	}

	// 从 URL 提取文件路径
	// URL 格式：/files/images/20260125_abc12345_image.jpg
	parts := strings.Split(request.URL, "/")
	if len(parts) < 3 {
		utils.HandleResponse(c, http.StatusBadRequest, nil, "Invalid file URL")
		return
	}

	fileType := parts[2] // images 或 audios
	filename := parts[len(parts)-1]

	// 构建实际文件路径
	filePath := filepath.Join("uploads", fileType, filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.HandleResponse(c, http.StatusNotFound, nil, "File not found")
		return
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, nil, "Failed to delete file")
		return
	}

	utils.HandleResponse(c, http.StatusOK, nil, "File deleted successfully")
}
