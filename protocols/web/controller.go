package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/services"
	"github.com/tim0-12432/simple-test-server/docker"
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

	// Upload a file and copy it into the nginx container webroot
	web.POST("/:id/upload", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		// Ensure this is a WEB server
		if strings.ToUpper(container.Type) != "WEB" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "container is not a web server"})
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		// Save to temp via web service
		res, err := SaveUploadedFileToTmp(ctx, fileHeader)
		if err != nil {
			switch err {
			case ErrMissingFile:
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
			case ErrInvalidType:
				c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			case ErrTooLarge:
				c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save uploaded file"})
			}
			return
		}

		defer func() { _ = os.Remove(res.LocalPath) }()

		// copy into container's web root using docker helper
		containerName := container.Name
		destPath := filepath.Join("/usr/share/nginx/html", res.SafeName)
		if err := docker.CopyFileToContainer(ctx, containerName, res.LocalPath, destPath, 30*time.Second); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to copy file to container: %v", err)})
			return
		}

		// Determine host port that maps to container port 80
		hostPort := 80
		if p, ok := container.Ports[80]; ok {
			hostPort = p
		}

		url := fmt.Sprintf("http://localhost:%d/%s", hostPort, res.SafeName)
		c.JSON(http.StatusCreated, gin.H{"url": url})
	})
}
