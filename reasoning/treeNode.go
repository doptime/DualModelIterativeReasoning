package reasoning

import (
	"DualModelIterativeReasoning/models"
	"strconv"
	"strings"

	"github.com/doptime/doptime/db"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var KeyTreeNode = db.HashKey[string, *TreeNode]()
var NodesMap = cmap.New[*TreeNode]()

type TreeNode struct {
	Id       string
	ParentId string
	Layer    int

	UserMsg      *models.Message
	Solution     *models.Message
	Verification *models.Message
	score        float64 `msgpack:"-"`
}

func (node *TreeNode) Score() (Score float64) {
	var err error
	if node.score != 0 {
		return node.score
	}
	if node.Verification == nil || len(node.Verification.Content) == 0 {
		return 0
	}
	s := strings.ToLower(node.Verification.Content)
	items := strings.Split(s, "score:")
	for _, txt := range items {
		txt = strings.TrimSpace(txt)
		//reserve number char or . only in leading text
		validInd := 0
		for ; validInd < len(txt) && strings.ContainsRune("0123456789.", rune(txt[validInd])); validInd++ {
		}
		if validInd == 0 {
			continue
		}
		txt = txt[:validInd]
		node.score, err = strconv.ParseFloat(txt, 64)
		if err == nil {
			break
		}
	}
	return node.score
}
func (node *TreeNode) Refinement(leadingtext string) (refinementMsg *models.Message) {
	if node.Verification == nil || len(node.Verification.Content) == 0 {
		return nil
	}
	leadingtext = strings.ToLower(leadingtext)
	s := strings.ToLower(node.Verification.Content)
	ind := strings.Index(s, leadingtext)
	if ind < 0 {
		return nil
	}
	text := strings.TrimSpace(node.Verification.Content[ind : len(node.Verification.Content)-1])
	text = strings.Split(text, "##")[0]
	return models.AssistantMsg(text)
}

func (parent *TreeNode) NewChild(id string) (newNode *TreeNode) {

	newNode = &TreeNode{Id: id, ParentId: parent.Id, Layer: parent.Layer + 1}
	NodesMap.Set(id, newNode)
	return newNode
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
func (node *TreeNode) BestScoreNode() (bestChild *TreeNode) {
	value := float64(0)
	NodesMap.IterCb(func(key string, node *TreeNode) {
		if node.Score() > value && node.Layer >= 1 {
			value = node.Score()
			bestChild = node
		}
	})
	return bestChild
}
