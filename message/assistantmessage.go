package message

func Assistant(msg string) *Message {
	return &Message{Role: "assistant", Content: msg}
}
