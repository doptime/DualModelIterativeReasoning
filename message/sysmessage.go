package message

func SysMsg(msg string) *Message {
	return &Message{Role: "system", Content: msg}
}
