package message

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (msg *Message) String() string {
	if msg == nil {
		return ""
	}
	return msg.Role + ": " + msg.Content
}

func SysMsg(msg string) *Message {
	if msg == "" {
		return nil
	}
	return &Message{Role: "system", Content: msg}
}
func UserMsg(msg string) *Message {
	if msg == "" {
		return nil
	}
	return &Message{Role: "user", Content: msg}
}
