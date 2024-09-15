package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/doptime/doptime/db"
	cmap "github.com/orcaman/concurrent-map/v2"
)

var KeyTreeNode = db.HashKey[string, *TreeNode]()
var NodesMap = cmap.New[*TreeNode]()

type TreeNode struct {
	Id         string
	ParentId   string
	Layer      int
	Difficulty float64

	UserMsg      *models.Message
	Solution     *models.Message
	Verification *models.Message
	score        float64 `msgpack:"-"`
}

var GetID = func() func(layter int) string {
	var IDHeader string = db.NanoId(4)
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

func ReadFloatAfterTag(tag, s string) (float64, error) {
	ind := strings.Index(s, tag)
	if ind < 0 {
		return 0, nil
	}
	s = s[ind+len(tag):]
	s = strings.TrimSpace(s)

	s = strings.Split(s, "\n")[0]
	if ind := strings.Index(s, "="); ind >= 0 {
		s = s[ind+1:]
	}
	s = strings.TrimSpace(s)

	validInd := 0
	for ; validInd < len(s) && strings.ContainsRune("0123456789.", rune(s[validInd])); validInd++ {
	}
	s = s[:validInd]
	if validInd == 0 {
		return 0, fmt.Errorf("no number found")
	}
	return strconv.ParseFloat(s, 64)
}

func (node *TreeNode) Score() (Score float64, err error) {
	if node.score != 0 {
		return node.score, nil
	}
	if node.Verification == nil || len(node.Verification.Content) == 0 {
		return 0, fmt.Errorf("no verification found")
	}
	s := strings.ToLower(node.Verification.Content)
	node.score, err = ReadFloatAfterTag("overall score:", s)
	if err != nil || node.score == 0 {
		node.score, err = ReadFloatAfterTag("evaluation:", s)
	}

	return node.score, err
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
	text = strings.Split(text, "**")[0]

	return models.AssistantMsg(text)
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
		score, err := node.Score()
		if node.Layer > 1 {
			score = score + math.Log10(float64(node.Layer))
		}
		if err != nil {
			return
		}
		if score > value && node.Layer >= 1 {
			value = score
			bestChild = node
		}
		if bestChild != nil && score == value && node.Layer >= bestChild.Layer {
			bestChild = node
		}
	})
	return bestChild
}
