package servers

import (
	"fmt"
)

type ServerDefinition interface {
	GetImage() string
	GetName() string
	GetPorts() []int
	GetEnv() map[string]string
}

type ServerInformation struct {
	Name  string            `json:"name"`
	Image string            `json:"image"`
	Ports []int             `json:"ports"`
	Env   map[string]string `json:"env"`
}

func GetAllServers() []ServerInformation {
	servers := []ServerDefinition{
		MqttServer{},
		WebServer{},
		FtpServer{},
		SmbServer{},
		MailServer{},
		OtelServer{},
	}
	var serverInfo []ServerInformation
	for _, server := range servers {
		info := ServerInformation{
			Name:  server.GetName(),
			Image: server.GetImage(),
			Ports: server.GetPorts(),
			Env:   server.GetEnv(),
		}
		serverInfo = append(serverInfo, info)
	}
	return serverInfo
}

func GetServerByType(serverType string) (*ServerInformation, error) {
	var serverDefinition ServerDefinition
	switch serverType {
	case "MQTT":
		serverDefinition = MqttServer{}
	case "WEB":
		serverDefinition = WebServer{}
	case "FTP":
		serverDefinition = FtpServer{}
	case "SMB":
		serverDefinition = SmbServer{}
	case "MAIL":
		serverDefinition = MailServer{}
	case "OTEL":
		serverDefinition = OtelServer{}
	default:
		return nil, fmt.Errorf("unknown server type: %s", serverType)
	}
	if serverDefinition == nil {
		return nil, fmt.Errorf("server type %s not found", serverType)
	}
	return &ServerInformation{
		Name:  serverDefinition.GetName(),
		Image: serverDefinition.GetImage(),
		Ports: serverDefinition.GetPorts(),
		Env:   serverDefinition.GetEnv(),
	}, nil
}
