package servers

type MailServer struct{}

func (s MailServer) GetImage() string {
	return "mailhog/mailhog:latest"
}

func (s MailServer) GetName() string {
	return "mail"
}

func (s MailServer) GetPorts() []int {
	return []int{1025, 8025}
}

func (s MailServer) GetEnv() map[string]string {
	return map[string]string{}
}
