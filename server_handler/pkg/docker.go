package pkg

import (
	"context"
	"fmt"
	"io"

	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func WriteToContainer(ctx context.Context, containerName string, input string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	// List containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	// Find container ID by name
	var containerID string
	for _, c := range containers {
		if slices.Contains(c.Names, "/"+containerName) {
			containerID = c.ID
			break
		}
	}
	if containerID == "" {
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
		return err
	}
	defer hijackedResp.Close()

	// Write input to container
	input = input + "\n"
	_, err = io.WriteString(hijackedResp.Conn, input)
	if err != nil {
		return err
	}

	// Optionally close stdin write
	if closer, ok := hijackedResp.Conn.(interface{ CloseWrite() error }); ok {
		if err := closer.CloseWrite(); err != nil {
			return err
		}
	}

	return nil
}
