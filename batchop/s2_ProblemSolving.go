package batchop

import (
	"DualModelIterativeReasoning/message"
	"DualModelIterativeReasoning/query"
	"DualModelIterativeReasoning/tools"
)

func generateProblemSolvingPrompt(node *query.TreeNode, formerSolutionGenerated *message.Message) (msg *message.Message) {
	prompt := "You are a world-class powerfull AI reasoning agent, cooperative, innovative, carefull, reflective and helpfull. Together with your AI counterpart, you are solving problems through structured collaboration.;"
	prompt += "Problem Reformulated:\n" + node.UserMsg.Content + "\n\n"
	if formerSolutionGenerated != nil && formerSolutionGenerated.Content != "" {
		s := formerSolutionGenerated.Content
		if part := tools.ReadMarkdownTagOut(s, "Solution Generated"); part != "" {
			s = part
		}
		prompt += "Previous Solution Generated:\n" + s + "\n\n"
	}
	prompt += `Your goal is to solve the Reformulated Problem Step by Step, according to following steps:
	## Step 1: reasoing to make revisions to the previous step-by-stey solutions (if applicable) 
	- Evaluate the weaknesses of the solution plan in previous solution step.
	- reasoing to Add or remove steps in the solution plan.
	- reasoning to unfold one step further in the solution step.
	- reasoning to fold one step back in the solution step.
	- reasoning to reasnwer the subquestion in the solution step.
	- reasoning to rephrase the question or subquestion in the solution step.
	- Throughout the process, please pay attention to:
		- Keep objective and fair evaluation of each solution
		- Actively look for synergies between different solutions
		- Continue to focus on the core of the problem and avoid deviating from the topic
		- Be bold and innovative while retaining the advantages of the original solution
		`
	prompt += "\n## Step 2: ** Final step-by-step Solution Generated ** :\n"
	prompt += "Based on the above analysis, write out the full step-by-step (Chain of Thought) solution plan for the problem.\n"
	prompt += "Iteration to improve previous solutions if applicable\n"
	prompt += "Before unfolding a step. explain the plan of the step in the leading sentence. (e.g. \"Step 1: First, I will ...\", \"Step 2: Next, I will ...\")"
	return message.UserMsg(prompt)
}

func ParallelProblemSolving(node *query.TreeNode) (msg *query.TreeNode, err error) {
	//Step 1: generate solutions
	problemToSolve := node.NewChild("ProblemSolveIter1")
	problemToSolve.UserMsg = generateProblemSolvingPrompt(node, nil)
	ProblemIter1 := problemToSolve.CloneN(4)
	err = query.AskLLMParallelly(ProblemIter1...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(ProblemIter1...)

	ProblemIter2 := node.NewChild("ProblemSolveIter2").CloneN(4)
	for i := 0; i < 4; i++ {
		ProblemIter2[i].UserMsg = generateProblemSolvingPrompt(node, ProblemIter1[i].AssistantMsg)
	}
	err = query.AskLLMParallelly(ProblemIter2...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(ProblemIter2...)

	//Step3: choose the best problem reformulatied
	msg, err = ParallelEvaluator(ProblemIter2...)
	CopyToClipboard(msg)
	msgBest := msg.Clone()
	if s := tools.ReadMarkdownTagOut(msgBest.AssistantMsg.Content, "Problem Reformulated"); len(s) > 0 {
		msgBest.AssistantMsg.Content = s
	}
	msgBest.Save()
	return msgBest, err
}
