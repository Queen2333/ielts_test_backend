package main

import (
	// _ "ielts_test_backend/docs" // 导入自动生成的文档

	"github.com/Queen2333/ielts_test_backend/database"
	"github.com/Queen2333/ielts_test_backend/routes"
	"github.com/Queen2333/ielts_test_backend/utils"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// GetClientIP 获取用户IP
func GetClientIP(c *gin.Context) string {
	if c.Request.Header.Get("X-Forwarded-For") != "" {
		return c.Request.Header.Get("X-Forwarded-For")
	}
	return c.ClientIP()
}

func main() {

	// 创建Gin实例
	//r := gin.Default()
	//
	//// 定义路由
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "Hello, Go Gin!",
	//		"ip":      GetClientIP(c),
	//	})
	//})

	str := utils.GenerateRandomString(32)

	println(str)

	// 注册路由
	r := routes.SetupRouter()
	// 使用 Swagger UI 中间件
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化数据库连接
	err := database.InitializeDB("ielts_alex:Yx180236@tcp(172.25.138.133:3306)/ielts_database")
	if err != nil {
		// 处理连接错误
		panic(err)
	}
	defer database.GetDB().Close()

	// 启动Gin服务
	r.Run(":8081")
}
