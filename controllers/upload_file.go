package controllers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// @Summary 上传文件
// @Description 上传一个文件到服务器
// @Tags File
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "上传的文件"
// @Success 200 {object} models.ResponseData{data=models.UploadResponse}
// @Failure 400 {object} models.ResponseData{data=nil}
// @Failure 500 {object} models.ResponseData{data=nil}
// @Router /upload [post]
func UploadFile(c *gin.Context) {
	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not open file")
		return
	}
	defer src.Close()

	// 创建一个缓冲区和 multipart writer
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// 创建一个 form-data field 来上传文件
	fw, err := w.CreateFormFile("file", filepath.Base(file.Filename))
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not create form file")
		return
	}

	// 将文件内容复制到 multipart writer
	if _, err = io.Copy(fw, src); err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not copy file content")
		return
	}
	w.Close()

	// 创建一个 HTTP 请求并设置 multipart content type
	flaskServerURL := "http://localhost:5001/upload" // Flask 服务器的 URL
	req, err := http.NewRequest("POST", flaskServerURL, &b)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not create request")
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not send request to Flask server")
		return
	}
	defer resp.Body.Close()

	// 读取 Flask 服务器的响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Could not read Flask server response")
		return
	}

	// 返回 Flask 服务器的响应
	c.Data(resp.StatusCode, "application/json", body)

	// // 定义上传文件的路径
	// filename := filepath.Base(file.Filename)
	// if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
	// 	utils.HandleResponse(c, http.StatusInternalServerError, "", "upload file err: %s")
	// 	return
	// }

	// baseURL := fmt.Sprintf("%s%s", c.Request.URL.Scheme, c.Request.Host)
	// fileURL := strings.Join([]string{baseURL, "uploads", filename}, "/")

	// utils.HandleResponse(c, http.StatusOK, fileURL, "Success")
}