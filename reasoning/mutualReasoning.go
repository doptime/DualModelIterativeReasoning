package reasoning

import (
	"DualModelIterativeReasoning/message"
	"DualModelIterativeReasoning/models"
	"context"
	"strings"

	"github.com/atotto/clipboard"
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

var SysPromptBasic = "You are a world-class powerfull AI system, cooperative, innovative, reflective and helpfull, capable of complex reasoning. Together with another AI model, you are solving problems through structured collaboration.;\n\n"

func Buildmessages(problem string, preSolution1 *message.Message, preSolution2 *message.Message) (msg *message.Message) {
	var prompt strings.Builder
	prompt.WriteString(SysPromptBasic)
	prompt.WriteString("#Problem:\n" + problem + "\n;")
	if preSolution1 != nil && len(preSolution1.Content) > 0 {
		prompt.WriteString("\n;Previous solution 1:\n" + preSolution1.Content + "\n;")
	}
	if preSolution2 != nil && len(preSolution2.Content) > 0 {
		prompt.WriteString("\n;Previous solution 2:\n" + preSolution2.Content + "\n;")
	}

	prompt.WriteString(`Please follow these steps below to perform dual-model iterative reasoning:

1. Evaluate Problem Reformulation in previous solutions and make revisions:
	- Holistic Problem Exploration revisions:
		- Analyze the given problem from multiple perspectives.
		- Identify potential underlying issues or broader contexts that may not be immediately apparent.
		- Consider various stakeholders and their potential concerns.
	- Intent Discovery:
		- Probe deeper into the possible motivations behind the question.
		- Identify any implicit assumptions or biases in the problem statement.
		- Consider how different framings of the problem might lead to different solutions.

2. Problem Reformulation:
	Based on 1) Problem Reformulation in previous solutions. 2) your holistic analysis revisions and intent discovery revisionsin step1. 
	- Gives the Context of the problem.
	- Give the Boundary Conditions of the problem.
	- Give the Constraints of the problem.
	- Give Reformulated Problem to capture its essence more accurately.

3. make revisions to the previous solutions plan or continue with iterating the previous solutions plan step by step:
   - Evaluate the weaknesses of the solution plan in previous solutions.
   - Is it deeply & very specificlly dived into the problem?

4. Give the new version of the solution plan (do not actually start solving yet, just make a plan):
 	- Based on the above analysis, write out the full step-by-step solution plan for the problem.
	- These improvements should absorb the advantages of the two solutions and try to make up for their shortcomings.

5. Evaluate previous solutions:
   - Based on the above analysis, Evaluate the weaknesses of the two previous solutions.
	- Compare the two solutions and find out their similarities and differences. Analyze the possible impact of these differences.
   - Is it deeply & very specificlly dived into the problem?
   - Based on the above analysis, propose 2-3 possible improvement directions. These improvements should absorb the advantages of the two solutions and try to make up for their shortcomings.
	
6. Generate a new solution:
	- Combine the advantages of the previous solution and the improvements you proposed, and follow the step-to-step solution plan to generate a new comprehensive solution.
	- Throughout the process, please pay attention to:
		- Keep objective and fair evaluation of each solution
		- Actively look for synergies between different solutions
		- Continue to focus on the core of the problem and avoid deviating from the topic
		- Be bold and innovative while retaining the advantages of the original solution

7. Tell wheather the problem is well and fully solved 
	- If the problem is not solved, the next action is continue, else stop.
Problem Solved: <yes> or <no>
`)
	return message.User(prompt.String())
}
func (node *TreeNode) ParalleBeamSearchUsingDualModelIterativeReasoning(Depty int) (err error) {
	var Models = []*models.Model{models.SLM1, models.SLM2}
	parent1, parent2 := node, node
	var LoopCnt = 0
	for ; ; LoopCnt++ {
		slm1, slm2 := Models[LoopCnt%len(Models)], Models[(LoopCnt+1)%len(Models)]
		childNode1, childNode2 := parent1.NewChild(), parent2.NewChild()
		msgpack := Buildmessages(node.UserMsg.Content, parent1.UserMsg, parent2.UserMsg)
		g, _ := errgroup.WithContext(context.Background())
		g.Go(func() (err error) {
			childNode1.Solution, err = slm1.AskLLM(0.7, false, msgpack)
			childNode1.Complete = strings.Contains(childNode1.Solution.Content, "Problem Solved: yes")
			childNode1.Save()
			return err
		})
		g.Go(func() (err error) {
			childNode2.Solution, err = slm2.AskLLM(0.7, false, msgpack)
			childNode2.Complete = strings.Contains(childNode2.Solution.Content, "Problem Solved: yes")
			childNode2.Save()
			return err
		})
		err = g.Wait()

		var stringBuilder strings.Builder
		stringBuilder.WriteString("\n\n# Solution1:\n")
		stringBuilder.WriteString(childNode1.Solution.Content)
		stringBuilder.WriteString("\n\n# Solution2:\n")
		stringBuilder.WriteString(childNode2.Solution.Content)
		clipboard.WriteAll(stringBuilder.String())

		if childNode1.Complete || childNode2.Complete || LoopCnt > Depty {
			break
		}
	}
	return err
}
