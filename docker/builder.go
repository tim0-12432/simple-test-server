package docker

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

const CUSTOM_IMAGE_PATH = "./custom_images/"

func CheckIfExistsLocally(imageName string) bool {
	args := []string{"image", "inspect", imageName}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Printf("Running Docker command: docker %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("docker run failed: %v - %s", err, strings.TrimSpace(string(out)))
		return false
	}

	if strings.Contains(string(out), "No such image") {
		return false
	}
	return true
}

func BuildCustomDockerImage(imageName string) error {
	if !CheckIfExistsLocally(imageName) {
		path := CUSTOM_IMAGE_PATH + strings.Split(imageName, ":")[0]
		dockerfilePath := path + "/Dockerfile"

		args := []string{"build", "-t", imageName, "-f", dockerfilePath, path}

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		log.Printf("Running Docker command: docker %s", strings.Join(args, " "))
		cmd := exec.CommandContext(ctx, "docker", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("docker build failed: %v - %s", err, strings.TrimSpace(string(out)))
		}
		log.Printf("Docker image %s built successfully", imageName)
	} else {
		log.Printf("Docker image %s already exists locally, skipping build", imageName)
	}
	return nil
}
