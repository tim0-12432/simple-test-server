package docker

import (
	"log"
	"net/http"
	"strings"

	"github.com/tim0-12432/simple-test-server/docker/servers"
)

type ServerConfiguration struct {
	name         string            `json:"name"`
	portMapping  map[int]int       `json:"ports"`
	envVariables map[string]string `json:"env"`
}

func StartServer(serverType string, config ServerConfiguration) int {
	var server servers.ServerDefinition
	switch serverType {
	case "MQTT":
		server = servers.FtpServer{}
	case "WEB":
		server = servers.WebServer{}
	case "FTP":
		server = servers.FtpServer{}
	case "SMB":
		server = servers.SmbServer{}
	case "MAIL":
		server = servers.MailServer{}
	default:
		log.Printf("Unknown server type: %s", serverType)
		return http.StatusBadRequest
	}

	if strings.Contains(server.GetImage(), "simple-test-server-custom-") {
		BuildCustomDockerImage(server.GetImage())
	}

	if err := RunContainer(config, serverType, server.GetImage(), server.GetName(), server.GetPorts(), server.GetEnv()); err != nil {
		log.Fatalf("Failed to start server %s: %v", serverType, err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
