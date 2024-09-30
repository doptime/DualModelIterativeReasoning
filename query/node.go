package query

import (
	"fmt"
	"time"

	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"
)

type Query struct {
	Created int64

	Group string

	Model string

	MsgSys       string
	MsgUser      string
	MsgAssistant string
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
func (node *Query) WithMsgSys(msg string) (old *Query) {
	node.MsgSys = msg
	return node
}
func (node *Query) WithMsgUser(msg string) (old *Query) {
	node.MsgUser = msg
	return node
}
func (node *Query) WithMsgAssistant(msg string) (old *Query) {
	node.MsgAssistant = msg
	return node
}
func (node *Query) WithModel(model string) *Query {
	node.Model = model
	return node
}

func (parent *Query) NewChildren(Stage string, userMsgs ...string) (newNode []*Query) {
	for _, msg := range userMsgs {
		newNode = append(newNode, parent.NewChild(Stage).WithMsgUser(msg))
	}
	return newNode
}
func (node *Query) Clone() (newNode *Query) {
	if node == nil {
		return nil
	}

	newNode = &Query{Created: getUniqId(), Group: node.Group, Model: node.Model}
	newNode.MsgSys = node.MsgSys
	newNode.MsgUser = node.MsgUser
	newNode.MsgAssistant = node.MsgAssistant

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

	node.MsgAssistant, err = model.AskLLM(0.7, false, message.SysMsg(node.MsgSys), message.UserMsg(node.MsgUser))
	return err

}
