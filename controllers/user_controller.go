package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/models"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
)

// RegisterUser 注册用户控制器
func RegisterUser(c *gin.Context) {
	// 获取请求中的表单数据
	email := c.PostForm("email")
	code := c.PostForm("code")

	// 检查邮箱是否为空
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱是必填字段"})
		return
	}

	// 检查验证码是否为空
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码是必填字段"})
		return
	}

	// 获取用户 IP 地址
	ip := utils.GetClientIP(c)

	fmt.Println(ip)

	// 这里可以添加更多的用户注册逻辑，比如将用户信息保存到数据库
	// 调用登录，直接下发token，需要同时处理redis里面的验证码（设置过期，或者删除）

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}

func GetUserInfo(c * gin.Context){
	var request struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	utils.InitRedis()
	user_info, err := utils.Get(request.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user in Redis"})
		return
	}

	var userInfo models.UserQuery
	err = json.Unmarshal([]byte(user_info), &userInfo)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userInfo})
}

// GetAllUser 获取所有用户
func GetAllUser(c *gin.Context) {

	db := database.GetDB()
	// 在控制器中执行数据库操作
	rows, err := db.Query("SELECT user_id, username, email, created_at, updated_at FROM users")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.UserQuery

	for rows.Next() {
		var user models.UserQuery
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取所有用户", "data": users})
}

// CreateUser 新建用户
func CreateUser(c *gin.Context) {

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否为空
	if requestBody.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名是必填字段"})
		return
	}

	// 检查邮箱是否为空
	if !utils.IsValidEmail(requestBody.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱是必填字段"})
		return
	}

	// 检查密码是否为空
	if requestBody.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码是必填字段"})
		return
	}

	db := database.GetDB()

	// 检查用户名是否已存在
	var existingUser models.UserQuery
	err := db.QueryRow("SELECT user_id FROM users WHERE username = ?", requestBody.Username).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check username"})
		return
	}

	// 检查邮箱是否已存在
	err = db.QueryRow("SELECT user_id FROM users WHERE email = ?", requestBody.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not check email"})
		return
	}

	// 在数据库中插入用户
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", requestBody.Username, requestBody.Email, requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add user"})
		return
	}

	// 这里可以添加更多的用户注册逻辑，比如将用户信息保存到数据库
	// 调用登录，直接下发token，需要同时处理redis里面的验证码（设置过期，或者删除）

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}
