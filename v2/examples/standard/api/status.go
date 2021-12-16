package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masonhubco/rebar/v2/examples/standard/models"
)

func Status(info models.Status) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, info.Snapshot())
	}
}
