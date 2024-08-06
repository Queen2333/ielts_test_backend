package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateRandomString generates a random string of the given length.
func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

// IsValidEmail 判断是否是合法的邮箱地址
func IsValidEmail(email string) bool {
	// 使用正则表达式来验证邮箱格式
	// 这里使用了一个简单的正则表达式，你可以根据需要使用更复杂的表达式
	// 这个表达式只检查了邮箱的基本格式，实际中可能需要更严格的验证
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	valid := regexp.MustCompile(emailRegex).MatchString(email)
	return valid
}

func HandleResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{"code": statusCode, "data": data, "message": message})
}

// StringToList 将字符串形式的数组转换为整数数组
func StringToList(str string) []int {
    // 去掉前后的中括号
    str = strings.Trim(str, "[]")

    // 通过逗号分割字符串
	re := regexp.MustCompile(`\s*,\s*`)
	strValues := re.Split(str, -1)

    // 初始化整数数组
    intValues := make([]int, 0, len(strValues))
    // 遍历字符串数组并转换为整数
    for _, strVal := range strValues {
		strVal = strings.TrimSpace(strVal)
		if intVal, err := strconv.Atoi(strVal); err == nil {
			intValues = append(intValues, intVal)
		} else {
			fmt.Println("Error converting:", strVal, err)
		}
	}
    return intValues
}
