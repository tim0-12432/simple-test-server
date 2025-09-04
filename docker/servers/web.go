package servers

type WebServer struct{}

func (s WebServer) GetImage() string {
	return "simple-test-server-custom-nginx:latest"
}

func (s WebServer) GetName() string {
	return "web"
}

func (s WebServer) GetPorts() []int {
	return []int{80}
}

func (s WebServer) GetEnv() map[string]string {
	return map[string]string{}
}
