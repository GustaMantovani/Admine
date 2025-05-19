package commandexecuter

import (
	"context"
	"fmt"
	"io"
	"server_handler/internal/config"

	"slices"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func WriteToContainer(input string) error {
	ctx := context.Background()

	// Cria o cliente Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return err
	}

	containerName := config.GetInstance().ComposeContainerName

	// Obtém o ID do container pelo nome
	var containerID string
	for _, container := range containers {
		if slices.Contains(container.Names, "/"+containerName) { // Os nomes incluem a barra inicial
			containerID = container.ID
		}
	}

	if containerID == "" {
		return fmt.Errorf("container '%s' not found", containerName)
	}

	fmt.Println(containerID)

	// Anexa ao stdin do container
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

	// Escreve no stdin
	input = input + "\n"
	_, err = io.WriteString(hijackedResp.Conn, input)
	if err != nil {
		return err
	}

	// Envia EOF para o stdin, se necessário
	if closer, ok := hijackedResp.Conn.(interface{ CloseWrite() error }); ok {
		err = closer.CloseWrite()
		if err != nil {
			return err
		}
	}

	return nil
}
