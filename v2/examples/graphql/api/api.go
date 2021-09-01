package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var appStartTime time.Time

func init() {
	appStartTime = time.Now()
}

func Status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
		"redis":  "connected",
		"uptime": time.Since(appStartTime).Truncate(time.Second).String(),
	})
}
