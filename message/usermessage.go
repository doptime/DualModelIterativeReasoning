package message

func UserMsg(msg string) *Message {
	return &Message{Role: "user", Content: msg}
}
