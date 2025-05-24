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
	log.Printf("Executing command in container: %s", input)
	ctx := context.Background()

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

	// Attach to container stdin
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

	log.Printf("Command '%s' executed successfully", input)
	return nil
}
