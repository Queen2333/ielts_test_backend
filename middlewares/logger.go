package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建日志记录器
		log := logrus.New()
		log.SetFormatter(&logrus.TextFormatter{})
		log.SetOutput(getLogFile())

		// 记录请求的接口和参数
		path := c.Request.URL.Path
		method := c.Request.Method
		params := getRequestParams(c)

		log.Infof("Request: %s %s?%s", method, path, params)

		// 执行请求处理
		c.Next()

		// 检查是否有错误发生
		errors := c.Errors.ByType(gin.ErrorTypeAny)
		if len(errors) > 0 {
			// 记录错误信息
			err := errors.Last().Err
			log.Errorf("Error: %v", err)
		}
	}
}

func getRequestParams(c *gin.Context) string {
	// 获取 URL 查询参数
	queryParams := c.Request.URL.Query()

	// 获取 JSON 请求体参数
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// 恢复请求体
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 解析 JSON 请求体参数
	var jsonParams map[string]interface{}
	err = json.Unmarshal(bodyBytes, &jsonParams)
	if err != nil {
		return queryParams.Encode()
	}

	// 合并 URL 查询参数和 JSON 请求体参数
	for key, value := range jsonParams {
		queryParams.Set(key, fmt.Sprint(value))
	}

	return queryParams.Encode()
}

func getLogFile() io.Writer {
	// 获取当前日期
	today := time.Now().Format("2006-01-02")

	// 创建日志文件夹
	logPath := "./logs/"
	err := os.MkdirAll(logPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return os.Stdout
	}

	// 创建日志文件
	logFileName := filepath.Join(logPath, today+".log")
	fmt.Println(logFileName, "logFileName")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to create log file: %v\n", err)
		return os.Stdout
	}

	return logFile
}
