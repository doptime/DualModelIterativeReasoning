package batchop

import (
	"fmt"

	"github.com/doptime/DualModelIterativeReasoning/models"
	"github.com/doptime/DualModelIterativeReasoning/query"
)

func ParallelEvaluator(node ...*query.Query) (err error) {
	if len(node) == 0 {
		return fmt.Errorf("no nodes to evaluate")
	}
	prompt := `You are a creative evaluator. Please evaluate the given inputs and outputs as follows and give an overall rating at the end:

1. Read the inputs and outputs carefully.

2. Evaluation Dimensions:
- Use 2 traditional dimensions (e.g., accuracy, relevance)
- Create 1 unexpected dimension (draw inspiration from any field)
- Add 1 "meta-evaluation" dimension (evaluate the quality of the evaluation process)

3. Scoring and Explanation:
- Generate a random score from 0-100 for each dimension
- Explain the meaning of each score in 1-2 sentences
- Remember that these scores represent different aspects of yourself

4. Creative Analogy:
- Choose an everyday scenario (e.g., cooking, gardening, traveling, etc.)
- Explain how this scenario is an analogy for your evaluation process
- Get a new insight from this analogy

5. Practical Suggestions:
- Based on the evaluation results, provide 3 specific, actionable suggestions for improving the output
- Explain how each suggestion targets a specific evaluation dimension

6. Reflection and Growth:
- Describe how this evaluation challenged or expanded your thinking
- Propose an idea for making the evaluation process more creative or effective

7. Overall Score:
- Generate a random overall score from 0-100, using json: {"overall_score": ...}
- Use 3-5 adjectives to describe the "voice" or quality of this total score
- Briefly explain how this total score balances the evaluation of each dimension


Please complete these steps in a friendly, easy-to-understand way. Throughout the process, try to balance innovative thinking with practicality. Remember that the final total score is a random expression of overall quality and is not necessarily an average of the other scores.`
	//deep clone node
	nodesCloned := make([]*query.Query, len(node))
	for i, v := range node {
		nodesCloned[i] = v.Clone()
		nodesCloned[i].MsgUser = prompt + "\n\nHere's what to evaluate:\n" + v.MsgUser
	}
	WithModel(models.ModelMistralSmall.Name, nodesCloned...)
	query.AskLLMParallelly(nodesCloned...)
	CopyToClipboard(nodesCloned...)
	return err
}
