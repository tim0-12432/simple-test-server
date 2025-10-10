package smb

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
	. "github.com/tim0-12432/simple-test-server/protocols/common"
	webpkg "github.com/tim0-12432/simple-test-server/protocols/web"
)

// InitializeSmbProtocolRoutes registers SMB-related HTTP routes.
func InitializeSmbProtocolRoutes(root *gin.RouterGroup) {
	smb := root.Group("/smb")

	smb.GET("/:id/", func(c *gin.Context) {
		serverID := c.Param("id")
		_, err := services.GetContainer(serverID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
	})

	// List file tree entries inside the container's smb share
	smb.GET("/:id/filetree", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		if strings.ToUpper(container.Type) != "SMB" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "container is not an smb server"})
			return
		}

		relPath := c.Query("path")
		if relPath == "" {
			relPath = ""
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 6*time.Second)
		defer cancel()

		entries, truncated, err := docker.ListSmbDir(ctx, container.Name, relPath, 1000)
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

	// Upload a file into the SMB share
	smb.POST("/:id/upload", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
			return
		}

		if strings.ToUpper(container.Type) != "SMB" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "container is not an smb server"})
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

		destRel := c.PostForm("path")
		if destRel == "" {
			destRel = res.SafeName
		}

		containerName := container.Name
		// ensure no leading slash
		destRel = strings.TrimPrefix(destRel, "/")
		if err := UploadFileToContainer(ctx, containerName, destRel, res.LocalPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to copy file to container: %v", err)})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"path": filepath.Join("/samba/share", destRel), "size": res.Size})
	})

	// Fetch container logs (tail)
	smb.GET("/:id/logs", func(c *gin.Context) {
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
