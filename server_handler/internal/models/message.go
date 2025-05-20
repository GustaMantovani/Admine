package models

import (
	"encoding/json"
)

type Message struct {
	Tags []string `json:"tags"`
	Msg  string   `json:"message"`
}

func NewMessage(msg string, tags []string) Message {
	return Message{
		Tags: tags,
		Msg:  msg,
	}
}

func (msg Message) ToString() string {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
