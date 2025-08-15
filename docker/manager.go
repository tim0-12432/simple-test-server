package docker

import (
	"log"

	"github.com/moby/moby/client"
)

var DockerClient *client.Client

func InitializeDockerClient() {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	DockerClient = client
	log.Println("Docker client initialized successfully")
}

func RunContainer(config ServerConfiguration, image string, name string, ports []int, env map[string]string) error {
	// Placeholder for running a Docker container
	// This function should implement the logic to run a Docker container with the specified parameters
	log.Printf("Running container: %s with name: %s", image, name)
	return nil // Replace with actual implementation
}

func StopAllContainers() error {
	// Placeholder for stopping all Docker containers
	// This function should implement the logic to stop all running Docker containers
	log.Println("Stopping all Docker containers")
	return nil // Replace with actual implementation
}

func CloseDockerClient() error {
	if DockerClient != nil {
		if err := DockerClient.Close(); err != nil {
			log.Printf("Error closing Docker client: %v", err)
			return err
		}
		log.Println("Docker client closed successfully")
	}
	return nil
}
