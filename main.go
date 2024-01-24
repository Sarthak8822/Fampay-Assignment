// main.go
package main

import (
	"Fampay_Backend_Assignment/api"
	"Fampay_Backend_Assignment/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/videos", api.GetVideosHandler)
	router.GET("/latest-videos", api.GetLatestVideosHandler)

	// Start continuous background fetching in a goroutine
	go service.FetchAndStoreVideos("official")

	router.Run(":8080")
}
