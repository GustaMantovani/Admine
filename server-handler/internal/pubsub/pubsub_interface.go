package pubsub

type PusSubInterface interface {
	ListenForMessages()
	SendMessage(string)
}
