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

	// 听力
	r.GET("/config/listening/list", controllers.ListeningList)
	r.GET("/config/listening/detail/:id", controllers.ListeningDetail)
	r.POST("/config/listening/add", controllers.AddListening)
	r.PUT("/config/listening/update", controllers.UpdateListening)
	r.DELETE("/config/listening/delete/:id", controllers.DeleteListening)

	r.GET("/config/listening-part/list", controllers.ListeningPartList)
	r.GET("/config/listening-part/detail/:id", controllers.ListeningPartDetail)
	r.POST("/config/listening-part/add", controllers.AddListeningPart)
	r.PUT("/config/listening-part/update", controllers.UpdateListeningPart)
	r.DELETE("/config/listening-part/delete/:id", controllers.DeleteListeningPart)

	r.POST("/upload/file", controllers.UploadFile)

	// 阅读
	r.GET("/config/reading/list", controllers.ReadingList)
	r.GET("/config/reading/detail/:id", controllers.ReadingDetail)
	r.POST("/config/reading/add", controllers.AddReading)
	r.PUT("/config/reading/update", controllers.UpdateReading)
	r.DELETE("/config/reading/delete/:id", controllers.DeleteReading)

	r.GET("/config/reading-part/list", controllers.ReadingPartList)
	r.GET("/config/reading-part/detail/:id", controllers.ReadingPartDetail)
	r.POST("/config/reading-part/add", controllers.AddReadingPart)
	r.PUT("/config/reading-part/update", controllers.UpdateReadingPart)
	r.DELETE("/config/reading-part/delete/:id", controllers.DeleteReadingPart)

	// 写作
	r.GET("/config/writing/list", controllers.WritingList)
	r.GET("/config/writing/detail/:id", controllers.WritingDetail)
	r.POST("/config/writing/add", controllers.AddWriting)
	r.PUT("/config/writing/update", controllers.UpdateWriting)
	r.DELETE("/config/writing/delete/:id", controllers.DeleteWriting)

	r.GET("/config/writing-part/list", controllers.WritingPartList)
	r.GET("/config/writing-part/detail/:id", controllers.WritingPartDetail)
	r.POST("/config/writing-part/add", controllers.AddWritingPart)
	r.PUT("/config/writing-part/update", controllers.UpdateWritingPart)
	r.DELETE("/config/writing-part/delete/:id", controllers.DeleteWritingPart)

	// 测试 套题
	r.GET("/config/testing/list", controllers.TestingList)
	r.GET("/config/testing/detail/:id", controllers.TestingDetail)
	r.POST("/config/testing/add", controllers.AddTesting)
	r.PUT("/config/testing/update", controllers.UpdateTesting)
	r.DELETE("/config/testing/delete/:id", controllers.DeleteTesting)

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
