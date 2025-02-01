package message

import (
	"encoding/json"
	"server/handler/internal/docker"
)

type Message struct {
	Tags []string `json:"tags"`
	Msg  string   `json:"message"`
}

func New(msg string, tags []string) Message {
	return Message{
		Tags: tags,
		Msg:  msg,
	}
}

func GetMessageInJsonString(status, containerName string) string {
	var m Message
	m.Tags = append(m.Tags, status)
	m.Msg = docker.GetZeroTierNodeID(containerName)

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)

	return jsonString
}
