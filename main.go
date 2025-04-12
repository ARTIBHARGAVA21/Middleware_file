package main

import (
	"fmt"
	"middleware/middle"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/protected", middle.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Protected endpoint",
		})

	})

	fmt.Println("Server is running on port 8000")
	router.Run("0.0.0.0:8000")
}
