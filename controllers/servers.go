package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tim0-12432/simple-test-server/docker"
	"github.com/tim0-12432/simple-test-server/docker/servers"
)

func InitializeServerRoutes(root *gin.RouterGroup) {
	path := root.Group("/servers")

	path.GET("", func(c *gin.Context) {
		serverlist := servers.GetAllServers()
		c.JSON(http.StatusOK, serverlist)
	})

	path.GET("/:type", func(c *gin.Context) {
		serverType := c.Param("type")
		server, err := servers.GetServerByType(serverType)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, server)
	})

	path.POST("/:type", func(c *gin.Context) {
		var configuration docker.ServerConfiguration
		if err := c.ShouldBindJSON(&configuration); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration"})
			return
		}
		serverType := c.Param("type")
		reqId := uuid.New().String()
		go func() {
			docker.StartServerWithProgress(reqId, serverType, configuration)
		}()

		c.JSON(http.StatusAccepted, gin.H{"reqId": reqId})
	})
}
