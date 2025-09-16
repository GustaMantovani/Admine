package models

import (
	"github.com/GustaMantovani/Admine/server_handler/internal"
	"encoding/json"
)

type AdmineMessage struct {
	Origin  string   `json:"origin"`
	Tags    []string `json:"tags"`
	Message string   `json:"message"`
}

func NewAdmineMessage(tags []string, message string) *AdmineMessage {
	ctx := internal.Get()
	return &AdmineMessage{
		Origin:  ctx.Config.App.SelfOriginName,
		Tags:    tags,
		Message: message,
	}
}

func (m *AdmineMessage) HasTag(tag string) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (msg AdmineMessage) ToString() string {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
