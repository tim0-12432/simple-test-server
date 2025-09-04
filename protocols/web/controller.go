package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/services"
)

func InitializeWebProtocolRoutes(root *gin.RouterGroup) {
	web := root.Group("/web")

	web.GET("/:id/", func(c *gin.Context) {
		serverID := c.Param("id")
		_, err := services.GetContainer(serverID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
	})
}
