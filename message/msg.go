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
