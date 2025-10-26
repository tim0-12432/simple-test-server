package servers

type OtelServer struct{}

func (s OtelServer) GetImage() string {
	return "simple-test-server-custom-otel:latest"
}

func (s OtelServer) GetName() string {
	return "otel"
}

func (s OtelServer) GetPorts() []int {
	return []int{4317, 4318, 8888, 8889}
}

func (s OtelServer) GetEnv() map[string]string {
	return map[string]string{}
}
