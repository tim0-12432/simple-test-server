package mqtt

import (
	"fmt"
	"net/http"

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

		var url = "localhost:" + fmt.Sprint(container.Ports[1883])

		subscribeToMqtt(url, func(message []byte) {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
		})
	})
}
