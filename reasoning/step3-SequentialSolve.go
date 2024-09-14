package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

func (node *TreeNode) SequentialSolve(TryNum int) (err error) {

	SolverMesseges := models.UserMsg(`Given a question, and solution together with Improvement Suggestions. Further improve the solution to solve the problem, by follow these steps:
1. Problem Analysis:
Analyze the given problem, solution and refinement ideas. Identify key components, constraints, and potential approaches.

2. Plan reproposed:
- Devise a step-by-step plan to solve the problem. (don't actually start solving yet, just make a plan)

3. Solution to the question:
- Use Chain of Thought reasoning to work through the plan and write the full solution within thinking.
## Solution: <...>

4. Solution Refinement:
Evaluate the strengths and weaknesses of Plan Proposal`)
	VerifierMessege := models.UserMsg(`Given a question and Problem Analysis, Plan Proposal, Solution and it's refinement proposal, Follow these steps:
1. Multi-dimensional Scoring:
For each dimension, provide a score between 0 and 10 (allowing for decimal points), where:
0-2 = Incorrect or significantly flawed
3-5 = Partially correct but needs major improvements
6-8 = Mostly correct with minor improvements needed
9-10 = Excellent, with little to no room for improvement

Dimensions to score:
a) Correctness: How accurate and error-free is the solution?
b) Completeness: How thoroughly does it address all aspects of the problem?
c) Clarity: How clear and easy to understand is the explanation?
d) Efficiency: How optimal is the approach in terms of time/resource usage?
e) Originality: How innovative or creative is the solution?

Format:
Plan Correctness: <score> / 10
Reasoning: <brief explanation>

Plan Completeness: <score> / 10
Reasoning: <brief explanation>

Plan Clarity: <score> / 10
Reasoning: <brief explanation>

Plan Efficiency: <score> / 10
Reasoning: <brief explanation>

Solution Correctness: <score> / 10
Reasoning: <brief explanation>

Solution Completeness: <score> / 10
Reasoning: <brief explanation>

Solution Clarity: <score> / 10
Reasoning: <brief explanation>

Solution Efficiency: <score> / 10
Reasoning: <brief explanation>

2. Overall Evaluation:
Provide a weighted average of the above scores (you may assign different weights to each dimension based on their importance for the specific problem).

Overall Score: <weighted average, realvalue> / 10

<---->

3. Improvement Suggestions:
- Provide 1-3 specific suggestions for how the solution plan could be improved.
- consider one the following actions if there's no appropriate suggestion:
	- propose one step thought
	- propose next subquestion along with answer
	- propose answer subquestion again
	- propose rephrase question or subquestion
`)

	for i := 0; i < TryNum; i++ {
		var bestAnswerNode *TreeNode = MCTSTrajectory.BestScoreNode()
		NewBestTrail := bestAnswerNode.NewChild()
		if NewBestTrail.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, bestAnswerNode.Solution, bestAnswerNode.Refinement("Suggestions:"), SolverMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if NewBestTrail.Verification, err = models.SLM2.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, NewBestTrail.Solution, VerifierMessege); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		NewBestTrail.Save()
	}
	return nil
}
