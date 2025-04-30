package message

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
