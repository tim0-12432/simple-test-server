package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/docker"
	"github.com/tim0-12432/simple-test-server/docker/servers"
)

func InitializeServerRoutes(root *gin.RouterGroup) {
	path := root.Group("/servers")

	path.GET("", func(c *gin.Context) {
		serverlist := servers.GetAllServers()
		c.JSON(http.StatusOK, serverlist)
	})

	path.POST("/:type", func(c *gin.Context) {
		var configuration docker.ServerConfiguration
		if err := c.ShouldBindJSON(&configuration); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid configuration"})
			return
		}
		serverType := c.Param("type")
		status := docker.StartServer(serverType, configuration)
		if status != http.StatusOK {
			c.JSON(status, gin.H{"error": "Failed to start server"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Server started successfully"})
	})
}
