package servers

type SmbServer struct{}

func (s SmbServer) GetImage() string {
	return "ghcr.io/servercontainers/samba:smbd-only-latest"
}

func (s SmbServer) GetName() string {
	return "smb"
}

func (s SmbServer) GetPorts() []int {
	return []int{139, 445}
}

func (s SmbServer) GetEnv() map[string]string {
	return map[string]string{}
}
