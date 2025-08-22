package servers

type MqttServer struct{}

func (s MqttServer) GetImage() string {
	return "simple-test-server-custom-mqtt:latest"
}

func (s MqttServer) GetName() string {
	return "mqtt"
}

func (s MqttServer) GetPorts() []int {
	return []int{1883, 9001}
}

func (s MqttServer) GetEnv() map[string]string {
	return map[string]string{
		"MQTT_USERNAME": "user",
		"MQTT_PASSWORD": "password",
	}
}
