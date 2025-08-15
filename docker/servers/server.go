package servers

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
