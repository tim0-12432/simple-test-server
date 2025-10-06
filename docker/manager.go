package docker

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/tim0-12432/simple-test-server/db/dtos"
	"github.com/tim0-12432/simple-test-server/db/services"
)

var DockerClient interface{}

var id = 0

func CheckAndPullImage(image string) error {
	log.Printf("Running Docker command: docker images -q %s", image)
	cmd := exec.Command("docker", "images", "-q", image)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker images failed: %v - %s", err, strings.TrimSpace(string(out)))
	}

	if strings.TrimSpace(string(out)) != "" {
		log.Printf("Docker image %s is already available locally", image)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	log.Printf("Running Docker command: docker pull %s", image)
	cmdPull := exec.CommandContext(ctx, "docker", "pull", image)
	pullOut, err := cmdPull.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker pull failed: %v - %s", err, strings.TrimSpace(string(pullOut)))
	}

	log.Printf("Successfully pulled Docker image %s: %s", image, strings.TrimSpace(string(pullOut)))
	return nil
}

func RunContainer(config ServerConfiguration, cType string, image string, name string, ports []int, env map[string]string) error {

	var allPorts = map[int]int{}
	var allEnv = map[string]string{}
	var finalName = "simple-test-server-" + name + "-" + fmt.Sprint(id)
	id++

	for _, p := range ports {
		allPorts[p] = p
	}
	for k, v := range env {
		allEnv[k] = v
	}

	// Merge ports from configuration (frontend provides an array of maps like [{"80":8080}])
	for _, mapping := range config.Ports {
		for hpStr, cp := range mapping {
			if hpStr == "" {
				continue
			}
			hp, err := strconv.Atoi(hpStr)
			if err != nil {
				// ignore malformed host port entries
				continue
			}
			allPorts[cp] = hp
		}
	}

	for k, v := range config.Env {
		allEnv[k] = v
	}
	if config.Name != "" {
		finalName = config.Name
	}

	args := []string{"run", "-d", "--name", finalName, "--label", "managed_by=simple-test-server"}

	for cp, hp := range allPorts {
		args = append(args, "-p", fmt.Sprintf("%d:%d", hp, cp))
	}
	for k, v := range allEnv {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, image)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Printf("Running Docker command: docker %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker run failed: %v - %s", err, strings.TrimSpace(string(out)))
	}

	services.CreateContainer(&dtos.Container{
		ID:          strings.TrimSpace(string(out)),
		Name:        finalName,
		Image:       image,
		CreatedAt:   time.Now().UnixMilli(),
		Environment: allEnv,
		Ports:       allPorts,
		Volumes:     map[string]string{},
		Networks:    []string{"host"},
		Status:      dtos.Running,
		Type:        cType,
	})

	log.Printf("Started container name=%s image=%s output=%s", name, image, strings.TrimSpace(string(out)))
	return nil
}

func StopAllContainers() error {
	log.Printf("Running Docker command: docker ps -aq -f label=managed_by=simple-test-server")
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

	log.Printf("Running Docker command: docker rm -f %s", strings.Join(ids, " "))
	args := append([]string{"rm", "-f"}, ids...)
	cmdRm := exec.Command("docker", args...)
	rmOut, err := cmdRm.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm failed: %v - %s", err, strings.TrimSpace(string(rmOut)))
	}

	for _, id := range ids {
		services.UpdateContainerStatus(id, dtos.Discarded)
	}

	log.Printf("Removed containers: %s", strings.TrimSpace(string(rmOut)))
	return nil
}

func StopContainer(containerId string) error {
	log.Printf("Running Docker command: docker rm -f %s", containerId)
	cmdRm := exec.Command("docker", "rm", "-f", containerId)
	rmOut, err := cmdRm.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker rm failed: %v - %s", err, strings.TrimSpace(string(rmOut)))
	}

	services.UpdateContainerStatus(containerId, dtos.Discarded)
	return nil
}
