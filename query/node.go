package query

import (
	"fmt"
	"time"

	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"
)

type Query struct {
	Created int64
	Group   string

	Model string

	MsgSys       *message.Message
	MsgUser      *message.Message
	MsgAssistant *message.Message
}
type QueryList []*Query

var getUniqId = func() func() int64 {
	var uniqTimeId int64 = 0
	return func() int64 {
		id := time.Now().Unix()
		if id <= uniqTimeId {
			id = uniqTimeId + 1
		}
		uniqTimeId = id
		return id
	}
}()

func (parent *Query) NewChild(Group string) (newNode *Query) {
	newNode = &Query{Created: getUniqId(), Group: Group, Model: parent.Model}
	return newNode
}
func (node *Query) WithMessage(msg *message.Message) (old *Query) {
	if msg.Role == "system" {
		node.MsgSys = msg
	} else if msg.Role == "user" {
		node.MsgUser = msg
	} else if msg.Role == "assistant" {
		node.MsgAssistant = msg
	}
	return node
}
func (node *Query) WithModel(model string) *Query {
	node.Model = model
	return node
}

func (parent *Query) NewChildren(Stage string, msgs ...*message.Message) (newNode []*Query) {
	for _, msg := range msgs {
		newNode = append(newNode, parent.NewChild(Stage).WithMessage(msg))
	}
	return newNode
}
func (node *Query) Clone() (newNode *Query) {
	if node == nil {
		return nil
	}

	newNode = &Query{Created: getUniqId(), Group: node.Group, Model: node.Model}
	if node.MsgSys != nil {
		newNode.MsgSys = message.SysMsg(node.MsgSys.Content)
	}
	if node.MsgUser != nil {
		newNode.MsgUser = message.UserMsg(node.MsgUser.Content)
	}
	if node.MsgAssistant != nil {
		newNode.MsgAssistant = message.Assistant(node.MsgAssistant.Content)
	}

	return newNode
}
func (node *Query) CloneN(n int) (newNode []*Query) {
	newNode = make([]*Query, n)
	for i := 0; i < n; i++ {
		newNode[i] = node.Clone()
	}
	return newNode
}

func (node *Query) Solute() (err error) {
	model, ok := models.Models[node.Model]
	if !ok {
		return fmt.Errorf("model not found")
	}
	node.MsgAssistant, err = model.AskLLM(0.7, false, node.MsgSys, node.MsgUser)
	return err

}
