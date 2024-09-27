package query

import (
	"DualModelIterativeReasoning/message"
	"DualModelIterativeReasoning/models"
	"fmt"
	"time"

	"github.com/doptime/doptime/db"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var KeyTreeNode = db.HashKey[string, *TreeNode]()
var NodesMap = cmap.New[*TreeNode]()

type TreeNode struct {
	Id     string
	RootId string
	Stage  string
	TimeAt int64

	Model string

	SysMsg       *message.Message
	UserMsg      *message.Message
	AssistantMsg *message.Message
	EvalScore    float64

	Complete bool
}

func (parent *TreeNode) NewChild(Stage string) (newNode *TreeNode) {
	id := db.NanoId(8)
	CreateAt := time.Now().Unix()
	newNode = &TreeNode{Id: id, RootId: parent.RootId, Stage: Stage, TimeAt: CreateAt, Model: parent.Model}
	NodesMap.Set(id, newNode)
	return newNode
}
func (parent *TreeNode) NewChildren(Stage string, msgs ...*message.Message) (newNode []*TreeNode) {
	for _, msg := range msgs {
		child := parent.NewChild(Stage)
		if msg.Role == "system" {
			child.SysMsg = msg
		} else if msg.Role == "user" {
			child.UserMsg = msg
		} else if msg.Role == "assistant" {
			child.AssistantMsg = msg
		}
		newNode = append(newNode, child)
	}
	return newNode
}
func (node *TreeNode) Clone() (newNode *TreeNode) {
	if node == nil {
		return nil
	}

	id := db.NanoId(8)
	newNode = &TreeNode{Id: id, RootId: node.RootId, Stage: node.Stage, Model: node.Model, Complete: node.Complete}
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
func (node *TreeNode) CloneN(n int) (newNode []*TreeNode) {
	newNode = make([]*TreeNode, n)
	for i := 0; i < n; i++ {
		newNode[i] = node.Clone()
	}
	return newNode
}

func (node *TreeNode) Solute() (err error) {
	model, ok := models.Models[node.Model]
	if !ok {
		return fmt.Errorf("model not found")
	}
	node.AssistantMsg, err = model.AskLLM(0.7, false, node.SysMsg, node.UserMsg)
	node.TimeAt = time.Now().Unix()
	return err

}

func (n *TreeNode) Save() {
	KeyTreeNode.HSet(n.Id, n)
}
