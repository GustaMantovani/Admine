package healthchecker

import (
	"fmt"
	"log"
	"net/http"
	"server/handler/internal/minecraft-server"
	"server/handler/internal/pubsub"
	"strings"
	"time"

	docker "github.com/fsouza/go-dockerclient"
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

func ContainerStatusEndpoint(serverName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client, err := docker.NewClientFromEnv()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err.Error())
		}

		fmt.Fprintf(w, "%s", SeeContainerStatus(client, serverName))

	}
}

func RunHealthCheckerWeb(ms minecraftserver.MinecraftServerContainerByCompose) {
	http.HandleFunc("/status", ContainerStatusEndpoint(ms.ContainerName))
	if err := http.ListenAndServe(":3132", nil); err != nil {
		log.Fatal("Erro ao inicializar serve: ", err)
	}
}

func sendAlertToServerHandler() {
	config := pubsub.GetConfigServerChannelFromDotEnv("HEALTH_CHECKER_CHANNEL")
	subscriber := pubsub.CreatePubsub(config.Addr, config.Channel)
	subscriber.SendMessage("down")
}

func healthCheckerAlert(ms minecraftserver.MinecraftServerContainerByCompose) {
	for {
		status := ms.SeeStatus()

		if !strings.Contains(status, "Up") || len(status) == 0 {
			log.Println("Enviando alerta. Status: ", status)
			sendAlertToServerHandler()
		}

		time.Sleep(time.Second * 1)
	}
}

func RunHealthChecker(ms minecraftserver.MinecraftServerContainerByCompose) {
	go healthCheckerAlert(ms)

	RunHealthCheckerWeb(ms)
}
