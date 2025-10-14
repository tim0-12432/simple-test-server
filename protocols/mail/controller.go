package mail

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/services"
)

func InitializeMailProtocolRoutes(root *gin.RouterGroup) {
	mail := root.Group("/mail")
	mail.GET("/:id/messages", listMessagesHandler)
	mail.GET("/:id/messages/:seq", getMessageHandler)
}

func listMessagesHandler(c *gin.Context) {
	id := c.Param("id")

	container, err := services.GetContainer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	fmt.Printf("Container ports: %+v\n", container.Ports) // Debugging line
	httpPort, ok := container.Ports[8025]
	if !ok || httpPort == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTTP port not found in container configuration"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	msgs, err := fetchEmailMessages(c.Request.Context(), "localhost", httpPort, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch messages: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": msgs})
}

func getMessageHandler(c *gin.Context) {
	id := c.Param("id")
	seqStr := c.Param("seq")

	seq, err := strconv.Atoi(seqStr)
	if err != nil || seq < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mail id"})
		return
	}

	container, err := services.GetContainer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	httpPort, ok := container.Ports[8025]
	if !ok || httpPort == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTTP port not found in container configuration"})
		return
	}

	msg, err := fetchSingleMessage(c.Request.Context(), "localhost", httpPort, seq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch messages: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}
