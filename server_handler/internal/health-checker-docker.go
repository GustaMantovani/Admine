package internal

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"log"
)

func SeeContainerStatus(client *docker.Client, containerName string) string {
	var containerStatus string

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == fmt.Sprintf("/%s", containerName) {
				containerStatus = container.Status
				break
			}
		}
	}

	return containerStatus
}
