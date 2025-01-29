package serverhandler

import (
	"log"
	"server/handler/internal/docker"
	"server/handler/internal/pubsub"
)

func RunServerHandler() {
	config := pubsub.GetConfigServerChannelFromDotEnv("HEALTH_CHECKER_CHANNEL")

	c := make(chan string)
	var msg string

	go pubsub.ListenChannelForMessages(config.Channel, config.Addr, c)

	for {
		msg = <-c
		if msg == "down" {
			log.Println("Subindo servidor")
			res, err := docker.StartServerDockerCompose("/home/andre/pgm/pessoal/Admine/minecraft-server/")
			if err != nil {
				log.Fatal(err)
			}

			r := string(res)

			log.Println(r)
		}
	}
}
