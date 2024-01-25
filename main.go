// main.go
package main

import (
	"Fampay_Backend_Assignment/api"
	"Fampay_Backend_Assignment/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Enable CORS (just in case if we have frontend)
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	// Routes
	router.GET("/api/paginated-videos", api.GetPaginatedVideosHandler)

	// Start continuous background fetching in a goroutine
	go func() {
		if err := service.FetchAndStoreVideos("official"); err != nil {
			log.Fatalf("Error in FetchAndStoreVideos: %v", err)
		}
	}()

	// Run the server
	router.Run(":8080")
}
