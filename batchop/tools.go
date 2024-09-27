package batchop

import (
	"DualModelIterativeReasoning/query"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func WithModel(Model string, node ...*query.TreeNode) {
	for _, v := range node {
		v.Model = Model
	}
}

func CopyToClipboard(node ...*query.TreeNode) {

	var stringBuilder strings.Builder
	for _, n := range node {
		timeAt := time.Unix(n.TimeAt, 0).Format("2006-01-02 15:04")
		stringBuilder.WriteString("\n\n# Stage: " + n.Stage + " Time: " + timeAt + " Model: " + n.Model + " Solution: \n\n")
		if n.Solution != nil {
			stringBuilder.WriteString(n.Solution.Content)
		}
	}

	if s := stringBuilder.String(); len(s) > 0 {
		clipboard.WriteAll(s)
	}
}
