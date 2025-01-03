package internal

import "encoding/json"

type Message struct {
	Tags []string `json:"tags"`
	Msg  string   `json:"message"`
}

func ConvertMessageToJson(status string) string {
	var m Message
	m.Tags = append(m.Tags, status)
	m.Msg = GetZeroTierNodeID()

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)

	return jsonString
}
