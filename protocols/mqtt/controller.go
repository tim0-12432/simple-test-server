package mqtt

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tim0-12432/simple-test-server/config"
	"github.com/tim0-12432/simple-test-server/db/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		var origin = r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://" + config.EnvConfig.Host + ":" + config.EnvConfig.Port,
		}
		if config.EnvConfig.AllowedOrigins != nil {
			allowedOrigins = append(allowedOrigins, config.EnvConfig.AllowedOrigins...)
		}
		if config.EnvConfig.Env == "DEV" {
			allowedOrigins = append(allowedOrigins, "http://localhost:5173")
		}
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	},
}

func InitializeMqttProtocolRoutes(root *gin.RouterGroup) {
	mqtt := root.Group("/mqtt")

	mqtt.GET("/:id/messages", func(c *gin.Context) {
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

		port := 1883
		if p, ok := container.Ports[port]; ok {
			port = p
		}
		var url = "localhost:" + fmt.Sprint(port)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// mutex to protect websocket writes
		var writeMutex sync.Mutex

		stop, err := startMqttSubscriber(ctx, url, func(message []byte) {
			writeMutex.Lock()
			defer writeMutex.Unlock()
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				// Log and cancel on write error
				log.Printf("websocket write error: %v", err)
				cancel()
			}
		})
		if err != nil {
			log.Printf("failed to start mqtt subscriber: %v", err)
			// close connection
			_ = conn.Close()
			return
		}
		defer stop()

		// reader to detect closure from client
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					// client likely closed connection
					log.Printf("websocket read error or closed: %v", err)
					cancel()
					return
				}
				// small sleep to avoid busy loop
				time.Sleep(10 * time.Millisecond)
			}
		}()

		// wait until context cancelled (either reader or write error)
		<-ctx.Done()
	})
}
