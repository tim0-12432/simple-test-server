package servers

type FtpServer struct{}

func (s FtpServer) GetImage() string {
	return "garethflowers/ftp-server:latest"
}

func (s FtpServer) GetName() string {
	return "ftp"
}

func (s FtpServer) GetPorts() []int {
	return []int{20, 21}
}

func (s FtpServer) GetEnv() map[string]string {
	return map[string]string{
		"FTP_USER": "user",
		"FTP_PASS": "password",
	}
}
