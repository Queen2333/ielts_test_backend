package routes

import (
	"net/http"

	"github.com/Queen2333/ielts_test_backend/controllers"
	"github.com/Queen2333/ielts_test_backend/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures the application's routes.
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Static("/uploads", "./uploads")
	
	r.Use(func(c *gin.Context) {
		if c.Request.URL.Path != "/login" && c.Request.URL.Path != "/send-code" {
			middlewares.JWTAuthMiddleware()(c)
		}
	})
	r.Use(enableCORS)
	// 使用 Logger 中间件
	r.Use(middlewares.Logger())

	r.POST("/login", controllers.LoginHandler)
	r.POST("/send-code", controllers.SendCodeHandler)
	// r.POST("/register", controllers.RegisterUser)
	r.GET("/user-info", controllers.GetUserInfo)

	r.GET("/config/listening/list", controllers.ListeningList)
	r.POST("/config/listening/add", controllers.AddListening)
	r.PUT("/config/listening/update", controllers.UpdateListening)
	r.DELETE("/config/listening/delete/:id", controllers.DeleteListening)

	r.GET("/config/listening-part/list", controllers.ListeningPartList)
	r.POST("/config/listening-part/add", controllers.AddListeningPart)
	r.PUT("/config/listening-part/update", controllers.UpdateListeningPart)

	r.POST("/upload/file", controllers.UploadFile)

	// r.GET("/users", controllers.GetAllUser)
	// r.POST("/create-user", controllers.CreateUser)

	// 执行Python脚本
	// r.GET("/execute-open-the-door", controllers.ExecutePythonScriptHandler)

	// Define routes and their handlers.

	// r.GET("/products/:id", middlewares.JWTAuthMiddleware(), controllers.GetProductByID)

	// 使用 Swagger UI 中间件
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
