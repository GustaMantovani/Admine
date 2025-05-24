package commandexecuter

import (
	"context"
	"fmt"
	"io"
	"log"
	"server_handler/internal/config"

	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func WriteToContainer(input string) error {
	ctx := context.Background()
	log.Printf("Writing command to container: %s", input)

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("Failed to create Docker client: %v", err)
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		log.Printf("Failed to list containers: %v", err)
		return err
	}

	containerName := config.GetInstance().ComposeContainerName
	log.Printf("Looking for container: %s", containerName)

	// Get container ID by name
	var containerID string
	for _, container := range containers {
		if slices.Contains(container.Names, "/"+containerName) { // Names include the leading slash
			containerID = container.ID
			break
		}
	}

	if containerID == "" {
		log.Printf("Container '%s' not found", containerName)
		return fmt.Errorf("container '%s' not found", containerName)
	}

	log.Printf("Found container with ID: %s", containerID)

	// Attach to container stdin
	log.Println("Attaching to container stdin...")
	hijackedResp, err := cli.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: false,
		Stderr: false,
	})
	if err != nil {
		log.Printf("Failed to attach to container: %v", err)
		return err
	}
	defer hijackedResp.Close()

	// Write to stdin
	input = input + "\n"
	log.Printf("Writing command to stdin: %s", input)
	_, err = io.WriteString(hijackedResp.Conn, input)
	if err != nil {
		log.Printf("Failed to write to stdin: %v", err)
		return err
	}

	// Send EOF to stdin if necessary
	if closer, ok := hijackedResp.Conn.(interface{ CloseWrite() error }); ok {
		err = closer.CloseWrite()
		if err != nil {
			log.Printf("Failed to close write: %v", err)
			return err
		}
	}

	log.Printf("Successfully executed command: %s", input)
	return nil
}
