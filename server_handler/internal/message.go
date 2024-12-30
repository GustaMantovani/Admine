package internal

type Message struct {
	Tags []string `json:"tags"`
	Msg  string   `json:"message"`
}
