package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lux208716/go-gin-project/database"
	"github.com/lux208716/go-gin-project/routes"
	"github.com/lux208716/go-gin-project/utils"
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

	// 初始化数据库连接
	err := database.InitializeDB("admin:hk6ic4!D26@tcp(10.98.163.112:3306)/xydb")
	if err != nil {
		// 处理连接错误
		panic(err)
	}
	defer database.GetDB().Close()

	// 启动Gin服务
	r.Run(":8081")
}
