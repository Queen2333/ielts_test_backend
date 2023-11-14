package routes

import (
	"github.com/lux208716/go-gin-project/controllers"
	"github.com/lux208716/go-gin-project/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the application's routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(enableCORS)
	// 使用 Logger 中间件
	r.Use(middlewares.Logger())

	r.POST("/login", controllers.Login)
	r.POST("/send-pin", controllers.SendPinController)
	r.POST("/register", controllers.RegisterUser)

	r.GET("/users", controllers.GetAllUser)
	r.POST("/create-user", controllers.CreateUser)

	// 执行Python脚本
	r.GET("/execute-open-the-door", controllers.ExecutePythonScriptHandler)

	// Define routes and their handlers.

	r.GET("/products/:id", middlewares.JWTAuthMiddleware(), controllers.GetProductByID)

	return r
}

// enableCORS is a middleware function to enable CORS headers.
func enableCORS(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Handle the OPTIONS preflight request
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	c.Next()
}
