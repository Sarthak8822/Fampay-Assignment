// api/video_handler.go
package api

import (
	"Fampay_Backend_Assignment/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const defaultPageSize = 10

func GetPaginatedVideosHandler(c *gin.Context) {
	// not mandatory and By-default it is 1
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	// not mandatory and By-default it is 10
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	//Service call to GetPaginatedVideos in service directory
	videos, err := service.GetPaginatedVideos(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, videos)
}
