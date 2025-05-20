package models

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("{tags:[%s], message:%s}", strings.Join(msg.Tags, ", "), msg.Msg)
}
