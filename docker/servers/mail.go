package servers

type MailServer struct{}

func (s MailServer) GetImage() string {
	return "servercontainers/mail-box:latest"
}

func (s MailServer) GetName() string {
	return "mail"
}

func (s MailServer) GetPorts() []int {
	return []int{25, 465, 587}
}

func (s MailServer) GetEnv() map[string]string {
	return map[string]string{}
}
