package main

import (
	"log"
	"server/handler/internal"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	minecraftServerContainerByCompose := internal.NewMinecraftServerContainerByCompose(client, "mine_server", "/home/andre/pgm/pessoal/minecraft-server-on-docker/")
	minecraftServerContainerByCompose.UpMinecraftServerContainerByCompose()

	for {
		minecraftServerContainerByCompose.VerifyContainerAndUpIfDown()
		time.Sleep(1 * time.Second)
	}
}
