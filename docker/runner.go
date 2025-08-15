package docker

import (
	"log"
	"net/http"

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
	case "mqtt":
		server = servers.MqttServer{}
	default:
		log.Printf("Unknown server type: %s", serverType)
		return http.StatusBadRequest
	}

	if err := RunContainer(config, server.GetImage(), server.GetName(), server.GetPorts(), server.GetEnv()); err != nil {
		log.Fatalf("Failed to start server %s: %v", serverType, err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
