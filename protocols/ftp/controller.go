package ftp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/db/services"
	"github.com/tim0-12432/simple-test-server/docker"
	webpkg "github.com/tim0-12432/simple-test-server/protocols/web"
	. "github.com/tim0-12432/simple-test-server/protocols/common"
)

// InitializeFtpProtocolRoutes registers FTP-related HTTP routes.
func InitializeFtpProtocolRoutes(root *gin.RouterGroup) {
	ftp := root.Group("/ftp")

	ftp.GET("/:id/", func(c *gin.Context) {
		serverID := c.Param("id")
		_, err := services.GetContainer(serverID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
	})

	// List file tree entries inside the container's ftproot
	ftp.GET("/:id/filetree", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		if strings.ToUpper(container.Type) != "FTP" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "container is not an ftp server"})
			return
		}

		relPath := c.Query("path")
		if relPath == "" {
			relPath = ""
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()

		entries, truncated, err := docker.ListFtpDir(ctx, container.Name, relPath, 1000)
		if err != nil {
			s := err.Error()
			if strings.Contains(s, "container not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
				return
			}
			if strings.Contains(s, "must be relative") || strings.Contains(s, "must not contain") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to list directory: %v", err)})
			return
		}

		out := make([]gin.H, 0, len(entries))
		for _, e := range entries {
			out = append(out, gin.H{
				"name":       e.Name,
				"path":       e.Path,
				"type":       e.Type,
				"size":       e.Size,
				"modifiedAt": e.ModifiedAt.Format(time.RFC3339),
			})
		}

		c.JSON(http.StatusOK, gin.H{"entries": out, "truncated": truncated})
	})

	// Upload a file into the FTP root
	ftp.POST("/:id/upload", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		if strings.ToUpper(container.Type) != "FTP" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "container is not an ftp server"})
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		res, err := webpkg.SaveUploadedFileToTmp(ctx, fileHeader)
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

		containerName := container.Name
		destPath := filepath.Join("/home/user", res.SafeName)
		if err := docker.CopyFileToContainer(ctx, containerName, res.LocalPath, destPath, 30*time.Second); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to copy file to container: %v", err)})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"path": destPath, "size": res.Size})
	})

	// Fetch container logs (tail)
	ftp.GET("/:id/logs", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		lines := 200
		if s := c.Query("lines"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 {
				lines = v
			}
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		logs, err := docker.GetContainerLogs(ctx, container.Name, lines)
		if err != nil {
			s := err.Error()
			if strings.Contains(s, "container not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get logs: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"logs": logs})
	})
}
