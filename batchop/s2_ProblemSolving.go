package batchop

import (
	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"
	"github.com/doptime/DualModelIterativeReasoning/query"
	"github.com/doptime/DualModelIterativeReasoning/tools"
)

func generateProblemSolvingPrompt(node *query.Query, formerSolutionGenerated *message.Message) (msg *message.Message) {
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

func ProblemSolving(node *query.Query) (msg []*query.Query, err error) {
	//Step 1: generate solutions
	UserMsg := generateProblemSolvingPrompt(node, nil)
	problemToSolve := node.NewChildren("SolutionIter", UserMsg, UserMsg, UserMsg, UserMsg)
	WithModel(models.ModelMistralSmall.Name, problemToSolve...)

	err = query.AskLLMParallelly(problemToSolve...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(problemToSolve...)

	ProblemIter2 := node.NewChild("SolutionIter").CloneN(4)
	for i := 0; i < 4; i++ {
		ProblemIter2[i].WithMessage(generateProblemSolvingPrompt(node, problemToSolve[i].AssistantMsg))
	}
	err = query.AskLLMParallelly(ProblemIter2...)
	if err != nil {
		return nil, err
	}
	CopyToClipboard(ProblemIter2...)

	//Step3: choose the best problem reformulatied
	err = ParallelEvaluator(ProblemIter2...)
	if err != nil {
		return nil, err
	}
	return ProblemIter2, err
}
