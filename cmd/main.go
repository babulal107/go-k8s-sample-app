package main

import (
	"fmt"
	"github.com/babulal107/go-k8s-sample-app/internal"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	// Load config from env variables
	config := internal.LoadConfig()

	log.Printf("configs object : %+v\n", config)

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

	log.Println("Starting server on port: ", config.AppPort)
	if err := r.Run(fmt.Sprintf(":%s", config.AppPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
