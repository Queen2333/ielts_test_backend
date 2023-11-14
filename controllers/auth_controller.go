package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lux208716/go-gin-project/utils"

	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure secret key

// Login handles user authentication and generates a JWT token.
func Login(c *gin.Context) {

	// 初始化Redis连接
	err := utils.InitRedis()

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to initialize Redis"})
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	// Replace this with your actual authentication logic.
	// For demonstration, we'll use hardcoded credentials.
	if username == "user" && password == "password" {
		// Create a new token.
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims (payload data) in the token.
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours.

		// Sign the token with the secret key and get the complete, signed token.
		signedToken, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}

		/*************************** ********************************/
		// 保存到redis
		utils.Set("mykey", signedToken, 0)

		user, err := utils.Get("mykey")

		fmt.Println(user, "test redis")

		// Return the token as a JSON response.
		c.JSON(http.StatusOK, gin.H{"token": signedToken, "user": user})
		return
	}

	// Authentication failed.
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}
