package query

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

func CopyToClipboard(title, model string, node ...*TreeNode) {
	var stringBuilder strings.Builder
	// write time of now
	Title := "\n\n# Time: " + time.Now().Format("2006-01-02 15:04:05") + " \n" + "\n\n# Title:" + title + " Model: " + model + " \n"
	fmt.Println(Title)
	stringBuilder.WriteString(Title)
	for i, childNode := range node {
		stringBuilder.WriteString("\n\n# Solution" + strconv.Itoa(i+1) + ":\n")
		stringBuilder.WriteString(childNode.Solution.Content)
	}
	clipboard.WriteAll(stringBuilder.String())
}
