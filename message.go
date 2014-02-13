package neural

type message struct {
	command 	string
	value    float64
	revision int64
}

func makeMessage(command string, value float64, revision int64) *message {
	message := new(message)
	message.command = command
	message.value = value
	message.revision = revision
	return message
}