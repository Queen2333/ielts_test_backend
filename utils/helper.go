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
	"github.com/davecgh/go-spew/spew"
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
	err := InitRedis()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Redis: %w", err)
	}

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

func ProcessPartList(c *gin.Context, results []map[string]interface{}, name string) ([]map[string]interface{}, error) {
	for i, result := range results {
		// 处理 []interface{} 类型的 part_list
		partListInterface, ok := result["part_list"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to parse part_list")
		}

		// 调用 GetPartDetails 获取详细信息
		details, err := GetPartDetails(partListInterface, name)
		if err != nil {
			return nil, err
		}

		// 将查询结果放回到对应的 part_list 中
		results[i]["part_list"] = details
	}
	return results, nil
}

func GetPartDetails(partListInterface []interface{}, name string) ([]map[string]interface{}, error) {
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
		partDetail, err := database.GetPartsByIds(name, []int{id})
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

func CalculateScore(partList []map[string]interface{}, submittedAnswers []models.AnswerItem) float64 {
	score := 0
	// 创建一个映射用于快速查找用户答案
	answerMap := make(map[string]interface{})
	for _, ans := range submittedAnswers {
		answerMap[ans.No] = ans.Answer
	}
	// spew.Dump(answerMap, partList, "answerMap")
	for _, part := range partList {
		// spew.Dump(part, "part2")
		typeList, ok := part["type_list"].([]interface{})
		if !ok {
			fmt.Printf("Skipping part due to type_list not being a valid type: %v (type: %T)\n", part["type_list"], part["type_list"])
			continue
		}

		for _, typeItemInterface := range typeList {
			typeItem, ok := typeItemInterface.(map[string]interface{}) // Type assertion
			if !ok {
				fmt.Println("Skipping typeItem due to invalid type:", typeItemInterface)
				continue
			}
			questionType, _ := typeItem["type"].(string)
			questionListInterface, ok := typeItem["question_list"].([]interface{})
			if !ok {
				spew.Dump(questionListInterface, "questionList skip")
				continue
			}
			for _, questionInterface := range questionListInterface {
				question, ok := questionInterface.(map[string]interface{})
				if !ok {
					fmt.Println("Skipping question due to invalid type:", questionInterface)
					continue
				}
				questionNo, ok := question["no"].(string)
				if !ok {
					fmt.Println("Skipping question due to 'no' not being a string:", question)
					continue
				}
				correctAnswer := question["answer"]

				// 检查用户是否提交了该题答案
				if answer, exists := answerMap[questionNo]; exists && correctAnswer != nil {
					score += calculateScoreByType(questionType, correctAnswer, answer)
					// fmt.Printf("Question No: %s, User Answer: %v, Correct Answer: %v, Points: %d\n", questionNo, answer, correctAnswer, score)
				}
			}
		}
	}
	spew.Dump(score, "total score")
	// 定义分数映射
	totalScore := score
	// 根据 score 的值返回对应的分数
	switch {
	case totalScore >= 39:
		return 9.0
	case totalScore >= 37:
		return 8.5
	case totalScore >= 35:
		return 8.0
	case totalScore >= 33:
		return 7.5
	case totalScore >= 30:
		return 7.0
	case totalScore >= 27:
		return 6.5
	case totalScore >= 23:
		return 6.0
	case totalScore >= 20:
		return 5.5
	case totalScore >= 16:
		return 5.0
	case totalScore >= 13:
		return 4.5
	case totalScore >= 10:
		return 4.0
	case totalScore >= 6:
		return 3.5
	case totalScore >= 4:
		return 3.0
	case totalScore == 3:
		return 2.5
	case totalScore == 2:
		return 2.0
	case totalScore == 1:
		return 1.0
	default:
		return 0
	}
}

// 根据题型计算得分
func calculateScoreByType(questionType string, correctAnswer, userAnswer interface{}) int {
	spew.Dump(questionType, correctAnswer, userAnswer, "questionType, correctAnswer, userAnswer")
	switch questionType {
	case "multi_choice":
		// 多选题处理
		if correctAns, ok := correctAnswer.([]interface{}); ok {
			if userAns, ok := userAnswer.([]interface{}); ok {
				return calculateMultiChoiceScore(correctAns, userAns)
			}
		}
	default:
		// 其他题型（单选、填空、匹配、地图等）
		if correctAnswer == userAnswer {
			return 1
		}
	}
	return 0
}

// 计算多选题得分
func calculateMultiChoiceScore(correctAnswers, userAnswers []interface{}) int {
	score := 0
	correctSet := make(map[interface{}]bool)

	// 将正确答案放入集合
	for _, correct := range correctAnswers {
		correctSet[correct] = true
	}

	// 检查用户答案是否在正确答案集合中
	for _, userAnswer := range userAnswers {
		if correctSet[userAnswer] {
			score++
		}
	}
	return score
}
