package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// create default gin router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		msg := fmt.Sprintf("Hello from Golang App! Path: %s", c.Request.URL.Path)
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": msg,
		})
	})

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hello World",
		})
	})

	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
