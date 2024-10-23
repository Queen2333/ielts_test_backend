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

	"github.com/Queen2333/ielts_test_backend/database"
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

func ProcessPartList(c *gin.Context, results []map[string]interface{}) ([]map[string]interface{}, error) {
	for i, result := range results {
		// 处理 []interface{} 类型的 part_list
		partListInterface, ok := result["part_list"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse part_list")
		}

		// 调用 GetPartDetails 获取详细信息
		details, err := GetPartDetails(partListInterface)
		if err != nil {
			return nil, err
		}

		// 将查询结果放回到对应的 part_list 中
		results[i]["part_list"] = details
	}
	return results, nil
}

func GetPartDetails(partListInterface []interface{}) ([]map[string]interface{}, error) {
	// 转换为字符串数组
	var partListStrArray []string
	for _, part := range partListInterface {
		partListStrArray = append(partListStrArray, fmt.Sprint(part))
	}

	partListStr := strings.Join(partListStrArray, ",")
	partList := StringToList(partListStr)

	// 查询 part_list 中的详细信息
	var details []map[string]interface{}
	for _, id := range partList {
		partDetail, err := database.GetPartsByIds("listening_part_list", []int{id})
		if err != nil {
			return nil, fmt.Errorf("failed to query listening parts: %w", err)
		}
		if len(partDetail) > 0 {
			details = append(details, partDetail[0])
		}
	}

	return details, nil
}

type Answer struct {
	No     string      `json:"no"`
	Answer interface{} `json:"answer"`
}

func CalculateScore(partList []map[string]interface{}, submittedAnswers []Answer) int {
	score := 0

	// 创建一个映射用于快速查找答案
	answerMap := make(map[string]interface{})
	for _, ans := range submittedAnswers {
		answerMap[ans.No] = ans.Answer
	}

	for _, part := range partList {
		for _, typeItem := range part["type_list"].([]interface{}) {
			questions := typeItem.(map[string]interface{})["question_list"].([]interface{})
			for _, question := range questions {
				q := question.(map[string]interface{})
				questionNo := q["no"].(string)

				if answer, exists := answerMap[questionNo]; exists {
					correctAnswer := q["answer"]
					if correctAnswer != nil {
						// 判断答案类型
						switch correctAnswer.(type) {
						case []interface{}: // 多选
							if userAnswers, ok := answer.([]interface{}); ok {
								for _, correct := range correctAnswer.([]interface{}) {
									for _, userAnswer := range userAnswers {
										if correct == userAnswer {
											score++
											break
										}
									}
								}
							}
						default: // 单选
							if answer == correctAnswer {
								score++
							}
						}
					}
				}
			}
		}
	}

	return score
}