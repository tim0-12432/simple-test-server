package mail

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/dtos"
	"github.com/tim0-12432/simple-test-server/db/services"
	"github.com/tim0-12432/simple-test-server/docker"
)

func InitializeMailProtocolRoutes(root *gin.RouterGroup) {
	mail := root.Group("/mail")
	mail.GET("/:id/messages", listMessagesHandler)
	mail.GET("/:id/messages/:seq", getMessageHandler)
	mail.GET("/:id/logs", getLogsHandler)
}

func listMessagesHandler(c *gin.Context) {
	id := c.Param("id")

	container, err := services.GetContainer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	httpPort, ok := container.Ports[MailHogWebPort]
	if !ok || httpPort == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTTP port not found in container configuration"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		n, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
		if n < 1 || n > MaxLimit {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("limit must be between 1 and %d", MaxLimit)})
			return
		}
		limit = n
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
	messageID := c.Param("seq")

	container, err := services.GetContainer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	httpPort, ok := container.Ports[MailHogWebPort]
	if !ok || httpPort == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "HTTP port not found in container configuration"})
		return
	}

	msg, err := fetchSingleMessage(c.Request.Context(), "localhost", httpPort, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch message: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func getLogsHandler(c *gin.Context) {
	serverID := c.Param("id")

	// Validate query params first (before checking container existence)
	tail := 500
	if t := c.Query("tail"); t != "" {
		if n, perr := strconv.Atoi(t); perr == nil {
			tail = n
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tail parameter"})
			return
		}
	}
	if tail < 1 || tail > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tail must be between 1 and 5000"})
		return
	}

	container, err := services.GetContainer(serverID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	if strings.ToUpper(container.Type) != "MAIL" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "container is not a mail server"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	fmt.Printf("INFO: %+v\n", container)

	// Fetch logs from Docker
	lines, truncated, err := docker.FetchContainerLogs(ctx, container.ID, tail, nil)
	if err != nil {
		if err == docker.ErrContainerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}
		if err == docker.ErrContainerNotRunning {
			// Return 409 with any lines and truncated flag
			out := make([]gin.H, 0, len(lines))
			for _, l := range lines {
				out = append(out, gin.H{"ts": l.TS.Format(time.RFC3339), "line": l.Line})
			}
			c.JSON(http.StatusConflict, gin.H{"error": "container not running", "lines": out, "truncated": truncated})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get logs: %v", err)})
		return
	}

	// Success: format and return logs
	out := make([]gin.H, 0, len(lines))
	for _, l := range lines {
		out = append(out, gin.H{"ts": l.TS.Format(time.RFC3339), "line": l.Line})
	}

	c.JSON(http.StatusOK, gin.H{"lines": out, "truncated": truncated, "container_running": container.Status == dtos.Running})
}
