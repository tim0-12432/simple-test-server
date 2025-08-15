package docker

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var DockerClient interface{}

func RunContainer(config ServerConfiguration, image string, name string, ports []int, env map[string]string) error {

	var allPorts = map[int]int{}
	var allEnv = map[string]string{}
	var finalName = "simple-test-server-" + name

	for _, p := range ports {
		allPorts[p] = p
	}
	for k, v := range env {
		allEnv[k] = v
	}
	for hp, cp := range config.portMapping {
		allPorts[hp] = cp
	}
	if config.name != "" {
		finalName = config.name
	}

	args := []string{"run", "-d", "--name", finalName, "--label", "managed_by=simple-test-server"}

	for hp, cp := range allPorts {
		args = append(args, "-p", fmt.Sprintf("%d:%d", hp, cp))
	}
	for k, v := range allEnv {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, image)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %v - %s", err, strings.TrimSpace(string(out)))
	}

	log.Printf("Started container name=%s image=%s output=%s", name, image, strings.TrimSpace(string(out)))
	return nil
}

func StopAllContainers() error {
	cmdList := exec.Command("docker", "ps", "-aq", "-f", "label=managed_by=simple-test-server")
	out, err := cmdList.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker ps failed: %v - %s", err, strings.TrimSpace(string(out)))
	}

	ids := strings.Fields(strings.TrimSpace(string(out)))
	if len(ids) == 0 {
		log.Println("No managed containers to stop")
		return nil
	}

	args := append([]string{"rm", "-f"}, ids...)
	cmdRm := exec.Command("docker", args...)
	rmOut, err := cmdRm.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm failed: %v - %s", err, strings.TrimSpace(string(rmOut)))
	}

	log.Printf("Removed containers: %s", strings.TrimSpace(string(rmOut)))
	return nil
}
