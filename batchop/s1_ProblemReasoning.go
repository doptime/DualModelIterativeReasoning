package batchop

import (
	"strings"

	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/query"
)

func ProblemReformulation(node *query.Query) (msg []*query.Query, err error) {
	var prompt strings.Builder
	prompt.WriteString("You are a world-class powerfull AI reasoning agent, cooperative, innovative, carefull, reflective and helpfull. Together with your AI counterpart, you are solving problems through structured collaboration.;")
	prompt.WriteString("#Problem:\n" + node.UserMsg.Content + "\n")
	prompt.WriteString(`In order to Reformulate the input problem statement into a clear, well-defined problem statement that  capture it's essence fully and accurately and suitable for solving with a language model. `)

	cs := node.NewChildren("ProblemReformulate",
		message.UserMsg(prompt.String()+`Your goal is to do the Holistic Problem Exploration:
	 - Analyze the given problem from multiple perspectives.
	 - Identify potential underlying issues or broader contexts that may not be immediately apparent.
	 - Consider various stakeholders and their potential concerns.`),
		message.UserMsg(prompt.String()+`Your goal is to do the Intent Discovery:
	 - Probe deeper into the possible motivations behind the question.
	 - Identify implicit assumptions or biases in the problem statement.
	 - Consider how different framings of the problem might lead to different solutions.`),
		message.UserMsg(prompt.String()+`Your goal is to do the Key Causual/context/Constraints Factors Identification:
	- List critical Causual factors that may influence the problem's solution.`),
		message.UserMsg(prompt.String()+`Your goal is to do the Key Result Identification:
	 - List critical factors Result may be introduced in the problem's solution.`),
	)
	err = query.AskLLMParallelly(cs...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(cs...)

	//ProblemReformulate
	prompt = strings.Builder{}
	prompt.WriteString("You are a world-class powerfull AI reasoning agent, cooperative, innovative, carefull, reflective and helpfull. Together with your AI counterpart, you are solving problems through structured collaboration.;")
	prompt.WriteString("#Problem:\n" + node.UserMsg.Content + "\n")
	prompt.WriteString(`given Problem Explorations on the Problem:`)
	prompt.WriteString("part1:\n" + cs[0].UserMsg.Content + "\n")
	prompt.WriteString("part2:\n" + cs[1].UserMsg.Content + "\n")
	prompt.WriteString("part3:\n" + cs[2].UserMsg.Content + "\n")
	prompt.WriteString("part4:\n" + cs[3].UserMsg.Content + "\n\n\n")
	prompt.WriteString(`### Final Reformulated Problem Statement
	Reformulate the input problem statement into a clear, well-defined problem statement that  capture it's essence fully and accurately according to following steps:
	- Provide the Context of the problem.
	- State the Condition/Constraints of the problem.
	- Present a reformulated problem statement (problem to solve) that captures its essence more accurately.
`)
	ProblemReformulate := node.NewChild("ProblemReformulate").WithMessage(message.UserMsg(prompt.String())).CloneN(4)
	err = query.AskLLMParallelly(ProblemReformulate...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(ProblemReformulate...)
	//choose the best problem reformulatied
	err = ParallelEvaluator(ProblemReformulate...)
	return ProblemReformulate, err
}
