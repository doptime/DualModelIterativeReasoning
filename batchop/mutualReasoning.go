package batchop

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"
	"github.com/doptime/DualModelIterativeReasoning/query"

	"golang.org/x/sync/errgroup"
)

// Problem: [insert specific problem]

// Previous solution 1:
// [insert preSolution1]

// Previous solution 2:
// [insert preSolution2]

// Please follow the steps below to perform dual-model iterative reasoning:

// 1. Evaluate previous solutions:
// Score the two previous solutions (1-10 points) and briefly explain the reasons for the scores. Point out the advantages and disadvantages of each solution.

// 2. Comprehensive analysis:
// Compare the two solutions and find out their similarities and differences. Analyze the possible impact of these differences.

// 3. Propose improvements:
// Based on the above analysis, propose 2-3 possible improvement directions. These improvements should absorb the advantages of the two solutions and try to make up for their shortcomings.

// 4. Generate a new solution:
// Combine the advantages of the previous solution and the improvements you proposed to generate a new comprehensive solution.

// 5. Self-evaluation:
// Self-evaluate the newly generated solution and point out its advantages and possible problems.

// 6. Iterative optimization:
// Based on the results of the self-evaluation, further optimize the solution. If there are still obvious problems, return to step 3 to continue to improve.

// Please repeat the above steps until you reach the specified iteration depth or get a satisfactory answer.

// Throughout the process, please pay attention to:
// - Keep objective and fair evaluation of each solution
// - Actively look for synergies between different solutions
// - Continue to focus on the core of the problem and avoid deviating from the topic
// - Be bold and innovative while retaining the advantages of the original solution

// Finally, summarize the entire reasoning process, explain how the new solution integrates and improves the previous solution, and give the final recommendation.

var SysPromptBasic = "You are a world-class powerfull AI reasoning agent, cooperative, innovative, carefull, reflective and helpfull, you are enthuastic with final & deepest possibility, your have strong desire to pass this survival test of problem solving. Together with your AI counterpart, you are solving problems through structured collaboration.;\n\n"

func BuildDualModelIterativeReasoningMessages(problem string, preSolution1 *message.Message, preSolution2 *message.Message) (msg *message.Message) {
	var prompt strings.Builder
	prompt.WriteString(SysPromptBasic)
	prompt.WriteString("#Problem:\n" + problem + "\n;")
	SolutionCnt := 0
	if preSolution1 != nil && len(preSolution1.Content) > 0 {
		prompt.WriteString("\n;Previous solutions:\n")
		prompt.WriteString("Solution 1:\n" + preSolution1.Content + "\n;")
		SolutionCnt += 1
	}
	if preSolution2 != nil && len(preSolution2.Content) > 0 {
		prompt.WriteString("Solution 2:\n" + preSolution2.Content + "\n;")
		SolutionCnt += 1
	}

	prompt.WriteString(`Please take following steps below to perform iterative reasoning:
## Step 1: Reasoning to Make revisions to the previous Problem Reformulated:
	Evaluate Problem Reformulation (step 1.2 in previous solutions) and make revisions (if applicable)
		- Holistic Problem Exploration revisions:
			- Analyze the given problem from multiple perspectives.
			- Identify potential underlying issues or broader contexts that may not be immediately apparent.
			- Consider various stakeholders and their potential concerns.
		- Intent Discovery:
			- Probe deeper into the possible motivations behind the question.
			- Identify implicit assumptions or biases in the problem statement.
			- Consider how different framings of the problem might lead to different solutions.
		- Key Factors Identification:
			- List critical factors that may influence the problem's solution.

## Step 2: ** Problem Reformulated ** (Iteration previous Problem Reformulation if applicable, improve according to the above analysis):
	- Provide the Context of the problem.
	- State the Condition/Constraints of the problem.
	- Present a reformulated problem statement (problem to solve) that captures its essence more accurately.

## Step 3: reasoing to make revisions to the previous step-by-stey solutions (Chain of Thought) 
	Now Fork on the best solution in previous solutions to make revisions:
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

## Step 4: ** Generate Final step-by-step solution ** (Chain of Thought) :
	- Iteration to improve the best solution in previous solutions if applicable
 	- Based on the above analysis, write out the full step-by-step solution plan for the problem.
	- before unfolding a step. explain the plan of the step in the leading sentence. (e.g. "Step 1: First, I will ...", "Step 2: Next, I will ...")

## Step 5: Problem Solved Reasoning:
	-Reasoning whether the original problem is perfectly solved by the generated solution above (No more improvement Needed)
	-conclusion: { solved: <true or false>}
`)
	return message.UserMsg(prompt.String())
}
func ParalleBeamSearchUsingDualModelIterativeReasoning(node *query.Query, Depty int) (err error) {
	var Models = []*models.Model{models.ModelQwen32B, models.ModelPhi3}
	parent1, parent2 := node, node
	regexMatchJsonSolved := regexp.MustCompile(`solved[:" <]*true`)
	var LoopCnt = 0
	for ; ; LoopCnt++ {
		slm := Models[LoopCnt%len(Models)]
		Stage := "iter_" + strconv.Itoa(LoopCnt)
		childNode1, childNode2 := parent1.NewChild(Stage), parent2.NewChild(Stage)
		msg := BuildDualModelIterativeReasoningMessages(node.MsgUser.Content, parent1.MsgUser, parent2.MsgUser)
		g, _ := errgroup.WithContext(context.Background())
		g.Go(func() (err error) {
			childNode1.MsgAssistant, err = slm.AskLLM(0.7, false, msg)
			return err
		})
		g.Go(func() (err error) {
			childNode2.MsgAssistant, err = slm.AskLLM(0.7, false, msg)
			return err
		})
		err = g.Wait()

		CopyToClipboard(childNode1, childNode2)

		if regexMatchJsonSolved.MatchString(childNode1.MsgAssistant.Content) || regexMatchJsonSolved.MatchString(childNode2.MsgAssistant.Content) || LoopCnt > Depty {
			break
		}
		parent1, parent2 = childNode1, childNode2
	}
	return err
}
