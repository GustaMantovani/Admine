package message

import "encoding/json"

type Message struct {
	Tags []string `json:"tags"`
	Msg  string   `json:"message"`
}

func ConvertMessageToJson(status, containerName string) string {
	var m Message
	m.Tags = append(m.Tags, status)
	// m.Msg = GetZeroTierNodeID(containerName)

	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	jsonString := string(jsonBytes)

	return jsonString
}
