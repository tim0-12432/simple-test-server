package otel

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tim0-12432/simple-test-server/config"
	"github.com/tim0-12432/simple-test-server/db/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		// allow empty origin (non-browser clients)
		if origin == "" {
			return true
		}
		// allow all origins in development
		if config.EnvConfig != nil && config.EnvConfig.Env == "DEV" {
			return true
		}
		allowedOrigins := []string{
			"http://" + config.EnvConfig.Host + ":" + config.EnvConfig.Port,
		}
		if config.EnvConfig.AllowedOrigins != nil {
			allowedOrigins = append(allowedOrigins, config.EnvConfig.AllowedOrigins...)
		}
		// allow localhost origins
		allowedOrigins = append(allowedOrigins, "http://localhost", "http://127.0.0.1")
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == origin {
				return true
			}
			if allowedOrigin == "http://localhost" && strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			if allowedOrigin == "http://127.0.0.1" && strings.HasPrefix(origin, "http://127.0.0.1") {
				return true
			}
		}
		return false
	},
}

func InitializeOtelProtocolRoutes(root *gin.RouterGroup) {
	otel := root.Group("/otel")

	// WebSocket endpoint for streaming OTEL collector logs
	otel.GET("/:id/logs", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// mutex to protect websocket writes
		var writeMutex sync.Mutex

		// start streaming container logs
		errChan := make(chan error, 1)
		go func() {
			errChan <- StreamOtelLogs(ctx, container.ID, func(line string) {
				writeMutex.Lock()
				defer writeMutex.Unlock()
				if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					log.Printf("websocket write error: %v", err)
					cancel()
				}
			})
		}()

		// reader goroutine to detect client closure
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					log.Printf("websocket read error or closed: %v", err)
					cancel()
					return
				}
			}
		}()

		// wait until context cancelled (either reader error, write error, or stream ended)
		select {
		case <-ctx.Done():
			// cancelled
		case err := <-errChan:
			if err != nil {
				log.Printf("log streaming error: %v", err)
			}
		}
	})

	// WebSocket endpoint for streaming OTEL telemetry (JSON-wrapped)
	otel.GET("/:id/telemetry", func(c *gin.Context) {
		serverID := c.Param("id")
		container, err := services.GetContainer(serverID)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var writeMutex sync.Mutex

		errChan := make(chan error, 1)
		go func() {
			errChan <- StreamOtelLogs(ctx, container.ID, func(line string) {
				writeMutex.Lock()
				defer writeMutex.Unlock()
				if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					log.Printf("websocket write error: %v", err)
					cancel()
				}
			})
		}()

		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					log.Printf("websocket read error or closed: %v", err)
					cancel()
					return
				}
			}
		}()

		select {
		case <-ctx.Done():
		case err := <-errChan:
			if err != nil {
				log.Printf("telemetry streaming error: %v", err)
			}
		}
	})
}
