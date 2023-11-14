package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lux208716/go-gin-project/models" // Update this with your actual project package name
	"github.com/lux208716/go-gin-project/utils"
)

// GetProductByID is the handler for the "/products/:id" endpoint.
func GetProductByID(c *gin.Context) {

	// 初始化Redis连接
	err := utils.InitRedis()

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to initialize Redis"})
		return
	}

	user, err := utils.Get("mykey")

	fmt.Println(user, "test redis")

	// Get the product ID from the request parameters.
	// 把参数转换成整型
	productID := c.Param("id")
	intVal, _ := strconv.Atoi(productID)

	// Convert the productID to an integer (you can use strconv.Atoi).
	// ...

	// Assuming products is a slice of Product as defined in models.
	// Find the product in the database or data store by ID.
	var product *models.Product
	for _, p := range models.Products {
		if p.ID == intVal {
			product = &p
			break
		}
	}

	// Check if the product exists.
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Return the product as JSON response.
	c.JSON(http.StatusOK, product)
}
