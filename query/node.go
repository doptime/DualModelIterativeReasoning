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
	Id       string
	ParentId string
	Layer    int
	Model    string

	SysMsg   *message.Message
	UserMsg  *message.Message
	Solution *message.Message
	Complete bool
}

var GetID = func() func(layter int) string {
	var IDHeader string = time.Now().Format("01-02-15-04")
	var GlobalID int = 0
	return func(layter int) string {
		GlobalID++
		return fmt.Sprintf("%s-%d-%d", IDHeader, layter, GlobalID)
	}
}()

func (parent *TreeNode) NewChild() (newNode *TreeNode) {
	id := GetID(parent.Layer + 1)
	newNode = &TreeNode{Id: id, ParentId: parent.Id, Layer: parent.Layer + 1}
	NodesMap.Set(id, newNode)
	return newNode
}
func (node *TreeNode) Clone() (newNode *TreeNode) {
	newNode = &TreeNode{Id: node.Id, ParentId: node.ParentId, Layer: node.Layer}
	return newNode
}
func (node *TreeNode) Solute() (err error) {
	model, ok := models.Models.Get(node.Model)
	if !ok {
		return fmt.Errorf("model not found")
	}
	if node.SysMsg != nil {
		node.Solution, err = model.AskLLM(0.7, false, node.SysMsg, node.UserMsg)
	} else {
		node.Solution, err = model.AskLLM(0.7, false, node.UserMsg)
	}
	return err

}

func (n *TreeNode) Save() {
	KeyTreeNode.HSet(n.Id, n)
}

func (node *TreeNode) Children() (children []*TreeNode) {
	NodesMap.IterCb(func(key string, child *TreeNode) {
		if child.ParentId == node.Id {
			children = append(children, child)
		}
	})
	return children

}
