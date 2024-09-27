package query

import (
	"fmt"

	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"

	"github.com/doptime/doptime/db"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var KeyTreeNode = db.HashKey[string, *Query]()
var NodesMap = cmap.New[*Query]()

type Query struct {
	Id     string
	RootId string
	Stage  string

	Model string

	SysMsg       *message.Message
	UserMsg      *message.Message
	AssistantMsg *message.Message

	EvalScore float64

	Complete bool
}
type QueryList []*Query

func (parent *Query) NewChild(Stage string) (newNode *Query) {
	id := db.NanoId(8)
	newNode = &Query{Id: id, RootId: parent.RootId, Stage: Stage, Model: parent.Model}
	NodesMap.Set(id, newNode)
	return newNode
}
func (node *Query) WithMessage(msg *message.Message) (old *Query) {
	if msg.Role == "system" {
		node.SysMsg = msg
	} else if msg.Role == "user" {
		node.UserMsg = msg
	} else if msg.Role == "assistant" {
		node.AssistantMsg = msg
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

	id := db.NanoId(8)
	newNode = &Query{Id: id, RootId: node.RootId, Stage: node.Stage, Model: node.Model, Complete: node.Complete}
	if node.SysMsg != nil {
		newNode.SysMsg = message.SysMsg(node.SysMsg.Content)
	}
	if node.UserMsg != nil {
		newNode.UserMsg = message.UserMsg(node.UserMsg.Content)
	}
	if node.AssistantMsg != nil {
		newNode.AssistantMsg = message.Assistant(node.AssistantMsg.Content)
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
	node.AssistantMsg, err = model.AskLLM(0.7, false, node.SysMsg, node.UserMsg)
	return err

}

func (n *Query) Save() {
	KeyTreeNode.HSet(n.Id, n)
}
