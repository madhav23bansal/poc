package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/madhav23bansal/poc/devops/loki-grafana-go/internal/log"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Initialize logger
	log.InitLogger()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		log.Info().Msg("Received ping request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Error().Err(err).Msg("Failed to start server")
	}
}
