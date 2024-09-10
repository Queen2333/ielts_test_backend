package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Queen2333/ielts_test_backend/models"
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

func ProcessRequest(c *gin.Context) (map[string]interface{}, error) {

	var request struct {
		Name   string 	`form:"name,omitempty"`  // form 标签表示从 URL 查询参数中获取
		Status *int   	`form:"status,omitempty"` 
		Type   *int   	`form:"type,omitempty"`
	}

	if err := c.ShouldBindQuery(&request); err != nil {
		fmt.Println("Error binding query:", err)
		HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return nil, err
	}

	fmt.Printf("Request: %+v\n", request)

	conditions := make(map[string]interface{})
	if request.Status != nil {
		conditions["status"] = *request.Status
	}
	if request.Type != nil {
		conditions["type"] = *request.Type
		if *request.Type == 3 {
			userID, err := GetUserIDFromToken(c)
			if err != nil {
				// 处理获取 user_id 失败的情况
				HandleResponse(c, http.StatusUnauthorized, "", err.Error())
				return nil, err
			}
			conditions["user_id"] = userID
		}
	}
	if request.Name != "" {
		conditions["name"] = "%" + request.Name + "%"
	}

	return conditions, nil
}
func GetUserIDFromToken(c *gin.Context) (string, error) {
	// 获取 Authorization 头中的 token
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 检查 token 是否为空
	if token == "" {
		return "", errors.New("missing token in Authorization header")
	}

	// 初始化 Redis 并获取 user 信息
	InitRedis()
	val, err := Get(token)
	if err != nil {
		return "", errors.New("failed to retrieve user from Redis")
	}

	// 检查 Redis 中是否有对应的用户数据
	if val == "" {
		return "", errors.New("no user data found for the provided token")
	}

	// 解析 user 信息
	var userInfo models.UserQuery
	err = json.Unmarshal([]byte(val), &userInfo)
	if err != nil {
		return "", errors.New("failed to parse user data from Redis")
	}

	// 检查解析后的 user 信息
	if userInfo.ID == "" {
		return "", errors.New("invalid user data: missing user ID")
	}

	// 返回 user_id
	return userInfo.ID, nil
}