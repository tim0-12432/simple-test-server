package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/dtos"
	"github.com/tim0-12432/simple-test-server/db/services"
	"github.com/tim0-12432/simple-test-server/docker"
)

func InitializeContainerRoutes(root *gin.RouterGroup) {
	path := root.Group("/containers")

	path.GET("", func(c *gin.Context) {
		result, err := services.ListRunningContainers()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	path.GET("/:id", func(c *gin.Context) {
		containerId := c.Param("id")
		result, err := services.GetContainer(containerId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, result)
	})

	path.DELETE("/:id", func(c *gin.Context) {
		containerId := c.Param("id")

		docker.StopContainer(containerId)

		services.UpdateContainerStatus(containerId, dtos.Discarded)

		c.Status(http.StatusNoContent)
	})
}
