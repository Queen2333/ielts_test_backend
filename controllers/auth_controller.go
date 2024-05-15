package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jwtSecret = []byte("qC2dACu+zgx94ALrfmCTESkxoqfCG4ItCWknbz+XmfTfWNDvFeuYKOGXgAKqSq+7Bdu8jXrxZfpWwE0K0jPLHw==") // Replace with a secure secret key

// @Summary 用户登录
// @Description 用户使用邮箱和验证码登录系统，如果用户不存在，则自动创建新用户
// @Accept json
// @Produce json
// @Param email body string true "用户邮箱地址"
// @Param code body string true "验证码"
// @Success 200 {object} models.ResponseData
// @Failure 400 {object} models.ResponseData
// @Failure 401 {object} models.ResponseData
// @Failure 500 {object} models.ResponseData
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	// 解析请求参数
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleResponse(c, http.StatusBadRequest, "", "Invalid request")
		return
	}

	// 获取验证码
	storedCode, err := utils.Get("verification_code_" + request.Email)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to retrieve verification code")
		return
	}

	// 验证码校验
	if request.Code != storedCode {
		utils.HandleResponse(c, http.StatusUnauthorized, "", "Invalid verification code")
		return
	}

	// 检查数据库中是否存在用户
	user, err := checkUserExists(request.Email)
	if err != nil {
		if database.IsNoRowsError(err) {
			// 符合条件的用户不存在
			user, err = createUser(request.Email)
			if err != nil {
				utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to create user")
				return
			}
		} else {
			// 其他错误
			utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to check user existence")
			return
		}
	}

	// 生成 JWT
	token, err := GenerateJWT(c, user.ID)
	if err != nil {
		utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to generate token")
		return
	}

	_, err = utils.Get(token)
	if err != nil {
		// 将 token 存储到 Redis 中，有效期为24小时
		userString, err := json.Marshal(user)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		err = utils.Set(token, userString, 24 * time.Hour)
		if err != nil {
			utils.HandleResponse(c, http.StatusInternalServerError, "", "Failed to store user in Redis")
			return
		}
	}

	// 返回 token
	utils.HandleResponse(c, http.StatusOK, gin.H{"token": token, "user": user}, "success")
}

// 检查数据库中是否存在指定邮箱的用户
func checkUserExists(email string) (models.UserQuery, error) {
	db := database.GetDB()
	// 实现检查用户是否存在的逻辑
	// 如果用户存在，返回 true，否则返回 false
	var user models.UserQuery
	query := "SELECT id, email, role_id FROM user_list WHERE email = ?"
    err := db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.RoleID)
	if err != nil {
        return models.UserQuery{}, err
    }


	return user, nil
}

// 在数据库中创建用户
func createUser(email string) (models.UserQuery, error) {
	// 生成 UUID
	userID := uuid.New().String()
	// 确保 userID 的长度不超过 36 个字符
	if len(userID) > 36 {
		userID = userID[:36]
	}
	// 实现创建用户的逻辑
	db := database.GetDB()
	_, err := db.Exec("INSERT INTO user_list (id, email, role_id) VALUES (?, ?, 0)", userID, email)
	if err != nil {
		return models.UserQuery{}, err
	}

	// 构造新用户信息
	newUser := models.UserQuery{
		ID: userID,
		Email: email,
		RoleID: 0,
	}

	return newUser, nil
}

// GenerateJWT 生成 JWT
func GenerateJWT(c *gin.Context, user_id string) (string, error) {
	tokenString, err := utils.Get("token_" + user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token in Redis"})
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["email"] = user_id
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // 设置过期时间为 24 小时
		// 签名 JWT
		tokenString, err = token.SignedString(jwtSecret)
		if err != nil {
			return "", err
		}
		err = utils.Set("token_" + user_id, tokenString, 24 * time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token in Redis"})
		}
	}
	
	return tokenString, nil
}