package protocols

import (
	"github.com/gin-gonic/gin"
	"github.com/tim0-12432/simple-test-server/protocols/ftp"
	"github.com/tim0-12432/simple-test-server/protocols/mail"
	"github.com/tim0-12432/simple-test-server/protocols/mqtt"
	"github.com/tim0-12432/simple-test-server/protocols/otel"
	"github.com/tim0-12432/simple-test-server/protocols/smb"
	"github.com/tim0-12432/simple-test-server/protocols/web"
)

func InitializeProtocolRoutes(root *gin.RouterGroup) {
	protocols := root.Group("/protocols")

	mqtt.InitializeMqttProtocolRoutes(protocols)
	web.InitializeWebProtocolRoutes(protocols)
	ftp.InitializeFtpProtocolRoutes(protocols)
	smb.InitializeSmbProtocolRoutes(protocols)
	mail.InitializeMailProtocolRoutes(protocols)
	otel.InitializeOtelProtocolRoutes(protocols)
}
