// api/video_handler.go
package api

import (
	"net/http"
	"strconv"

	"Fampay_Backend_Assignment/service"

	"github.com/gin-gonic/gin"
)

const defaultPageSize = 10

func GetVideosHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'query' is required"})
		return
	}

	// videos, err := service.FetchAndStoreVideos(query)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	// 	return
	// }

	// c.JSON(http.StatusOK, videos)
}

func GetLatestVideosHandler(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	// Fetch latest videos from MongoDB
	videos, err := service.GetLatestVideos(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, videos)
}
