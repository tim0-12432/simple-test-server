package controllers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/progress"
)

// registers SSE progress streaming
func InitializeProgressRoutes(root *gin.RouterGroup) {
	root.GET("/servers/progress/:reqId", func(c *gin.Context) {
		reqId := c.Param("reqId")
		ch, ok := progress.Default.Get(reqId)
		if !ok {
			ch = progress.Default.New(reqId)
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Flush()

		notify := c.Writer.CloseNotify()
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					// channel closed, send final event and return
					fmt.Fprintf(c.Writer, "data: {\"percent\":100,\"message\":\"done\",\"error\":false}\n\n")
					c.Writer.Flush()
					return
				}
				// send event as JSON
				fmt.Fprintf(c.Writer, "data: {\"percent\":%d,\"message\":\"%s\",\"error\":%v}\n\n", ev.Percent, ev.Message, ev.Error)
				c.Writer.Flush()
			case <-notify:
				return
			case <-time.After(30 * time.Second):
				// keepalive comment
				fmt.Fprintf(c.Writer, ": keepalive\n\n")
				c.Writer.Flush()
			}
		}
	})
}
