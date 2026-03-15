package pubsub

import "encoding/json"

// AdmineMessage is the message format used across all Admine pub/sub channels
type AdmineMessage struct {
	Origin  string   `json:"origin"`
	Tags    []string `json:"tags"`
	Message string   `json:"message"`
}

// NewAdmineMessage creates an AdmineMessage with an explicit origin
func NewAdmineMessage(origin string, tags []string, message string) *AdmineMessage {
	return &AdmineMessage{
		Origin:  origin,
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

func (m AdmineMessage) ToString() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
