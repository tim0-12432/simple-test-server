package docker

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tim0-12432/simple-test-server/docker/servers"
	"github.com/tim0-12432/simple-test-server/progress"
)

type ServerConfiguration struct {
	name         string            `json:"name"`
	portMapping  map[int]int       `json:"ports"`
	envVariables map[string]string `json:"env"`
}

func StartServerWithProgress(reqId string, serverType string, config ServerConfiguration) {
	progress.Default.New(reqId)
	progress.Default.Send(reqId, progress.Event{Percent: 10, Message: "starting", Error: false})

	var server servers.ServerDefinition
	switch serverType {
	case "MQTT":
		server = servers.MqttServer{}
	case "WEB":
		server = servers.WebServer{}
	case "FTP":
		server = servers.FtpServer{}
	case "SMB":
		server = servers.SmbServer{}
	case "MAIL":
		server = servers.MailServer{}
	default:
		msg := fmt.Sprintf("Unknown server type: %s", serverType)
		log.Printf(msg)
		progress.Default.Send(reqId, progress.Event{Percent: 100, Message: msg, Error: true})
	}

	if strings.Contains(server.GetImage(), "simple-test-server-custom-") {
		progress.Default.Send(reqId, progress.Event{Percent: 30, Message: "Building image", Error: false})
		if err := BuildCustomDockerImage(server.GetImage()); err != nil {
			progress.Default.Send(reqId, progress.Event{Percent: 50, Message: fmt.Sprintf("build failed: %v", err), Error: true})
			return
		}
		progress.Default.Send(reqId, progress.Event{Percent: 50, Message: "Build successful", Error: false})
	} else {
		progress.Default.Send(reqId, progress.Event{Percent: 30, Message: "Pulling image", Error: false})
		if err := CheckAndPullImage(server.GetImage()); err != nil {
			progress.Default.Send(reqId, progress.Event{Percent: 50, Message: fmt.Sprintf("pull failed: %v", err), Error: true})
			return
		}
		progress.Default.Send(reqId, progress.Event{Percent: 50, Message: "Pull successful", Error: false})
	}

	progress.Default.Send(reqId, progress.Event{Percent: 80, Message: "Starting container", Error: false})
	if err := RunContainer(config, serverType, server.GetImage(), server.GetName(), server.GetPorts(), server.GetEnv()); err != nil {
		progress.Default.Send(reqId, progress.Event{Percent: 90, Message: fmt.Sprintf("run failed: %v", err), Error: true})
		return
	}

	progress.Default.Send(reqId, progress.Event{Percent: 100, Message: "Started", Error: false})

	go func() {
		time.Sleep(30 * time.Second)
		progress.Default.Remove(reqId)
	}()
}
