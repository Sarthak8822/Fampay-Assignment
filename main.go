// main.go
package main

import (
	"Fampay_Backend_Assignment/api"
	"Fampay_Backend_Assignment/service"
	"log"

	"github.com/gin-gonic/gin"
)

// Use environment variables for configuration
const (
	apiRouteLatestVideos    = "/latest-videos"
	apiRoutePaginatedVideos = "/paginated-videos"
	defaultPort             = ":8080"
)

func main() {
	router := gin.Default()

	// Handle API errors with recovery middleware
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
	}))

	router.GET(apiRouteLatestVideos, api.GetLatestVideosHandler)
	router.GET(apiRoutePaginatedVideos, api.GetPaginatedVideosHandler)

	// Start continuous background fetching in a goroutine
	go func() {
		if err := service.FetchAndStoreVideos("official"); err != nil {
			log.Fatalf("Error in FetchAndStoreVideos: %v", err)
		}
	}()

	router.Run(defaultPort)
}
