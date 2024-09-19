package message

func User(msg string) *Message {
	return &Message{Role: "user", Content: msg}
}
