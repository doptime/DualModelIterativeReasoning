package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func (node *TreeNode) SequentialSolve(TryNum int) (err error) {

	SolverMesseges := models.MsgOfUser(`Given a question, and solution together with Improvement Suggestions. Further improve the solution to solve the problem, by follow these steps:

## 1. Problem Analysis and Intent Navigation
	- Holistic Problem Exploration:
		- Analyze the given problem from multiple perspectives.
		- Identify potential underlying issues or broader contexts that may not be immediately apparent.
		- Consider various stakeholders and their potential concerns.

	- Intent Discovery:
		- Probe deeper into the possible motivations behind the question.
		- Identify any implicit assumptions or biases in the problem statement.
		- Consider how different framings of the problem might lead to different solutions.

## 2. Question Reformulation:
	Based on your holistic analysis and intent discovery, reformulate the question to capture its essence more accurately.
	Provide a brief explanation of why this reformulation might lead to a more comprehensive or insightful answer.

## 3. Plan reproposed:
- Devise a step-by-step plan to solve the Question. (don't actually start solving yet, just make a plan)

## 4. Solving The Question Step by Step, According to the Plan:
Use Chain of Thought ,i.e., Work through the Plan Step by Step, write the full solution for each plan steps. until  finally answer the question.

## 5. Conclusion:

## 6. Evaluate Solution Refinement:
Evaluate the strengths and weaknesses of the solution to the question`)

	var bestAnswerNode *TreeNode = MCTSTrajectory.BestScoreNode()
	if bestAnswerNode == nil {
		fmt.Println("Error: bestAnswerNode is nil")
		return nil
	}
	for i := 0; i < TryNum; i++ {
		NewBestTrail := bestAnswerNode.NewChild()
		if NewBestTrail.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, bestAnswerNode.Solution, bestAnswerNode.Refinement("improvement Suggestions"), SolverMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if NewBestTrail.Verification, err = models.SLM2.AskLLM(0.7, false, MCTSTrajectory.UserMsg, NewBestTrail.Solution, models.MsgOfUser(`
Given a question, Problem Analysis, Plan Proposal, Solution, and its refinement proposal, follow these steps:

# Critical Analysis and Improvement Suggestions

## 1. Verify Problem Analysis and Intent Navigation
   - Evaluate the merits/weaknesses of Intent Navigation
   - What is dummb requirements in the problem analysis?
   - What is Redundant process in the problem analysis?

## 2. Verify Question Reformulated:
   - Evaluate the merits/weaknesses of Question Reformulated
   - Is it still fidelity to the original question?
   - Is it usefull & easy to take into action?	
   - Is it logically sound?

## 3. Verify the plan step by step:
   - Evaluate the weaknesses of the solution plan.
   - Is it deeply & very specificlly dived into the problem?

## 4. Verify the solution step by step:
   - Evaluate the weaknesses of the solution.

Methodology:
	1. Redundancy and Simplification:
	- Identify any redundant or unnecessary steps in the plan or solution.
	- Suggest how to simplify the approach while maintaining its effectiveness.

	2. Specific Improvements:
	- Provide specific, actionable suggestions to improve the solution plan.
	- For each suggestion, explain its potential impact on the overall solution.

	3. Alternative Perspective:
	- Propose the strongest possible argument or approach that could overturn the current conclusion.
	- Explain the reasoning behind this alternative perspective.

	4. Optional Refinements (choose 1-2 if applicable):
	- Propose a single-step plan to address a weakness in the current solution.
	- Suggest a sub-plan along with its potential answer to deepen the solution.
	- Recommend how to rephrase a step or sub-question for clarity or better focus.

Final Improvement Conclusion:
- Provide a concise summary of the most crucial improvements needed, emphasizing self-criticism and consideration of alternative perspectives.
		`)); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		// clear the console
		//copy to the clipboard
		var stringBuilder strings.Builder
		stringBuilder.WriteString("\n\n# Solution:\n")
		stringBuilder.WriteString(NewBestTrail.Solution.Content)
		stringBuilder.WriteString("\n\n# Verification:\n")
		stringBuilder.WriteString(NewBestTrail.Verification.Content)
		clipboard.WriteAll(stringBuilder.String())

		bestAnswerNode = NewBestTrail
		NewBestTrail.Save()
	}
	return nil
}
