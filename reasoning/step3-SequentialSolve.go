package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func (node *TreeNode) SequentialSolve(TryNum int) (err error) {

	SolverMesseges := models.UserMsg(`Given a question, and solution together with Improvement Suggestions. Further improve the solution to solve the problem, by follow these steps:
1. Problem Analysis:
Analyze the given problem, solution and refinement ideas. Identify key components, constraints, and potential approaches.

2. Plan reproposed:
- Devise a step-by-step plan to solve the problem. (don't actually start solving yet, just make a plan)

3. Solving The Question Step by Step, According to the Plan:
Use Chain of Thought ,i.e., Work through the Plan Step by Step, write the full solution for each plan steps. until  finally answer the question.

4. Conclusion:

5. Evaluate Solution Refinement:
Evaluate the strengths and weaknesses of the solution to the question`)

	VerifierMessege := models.UserMsg(`Given a question and Problem Analysis, Plan Proposal, Solution and it's refinement proposal, Follow these steps:
** Improvement Suggestions **
- reasoning to raise a most powerful plan to overturn the conclusion
- reasoning to remove redundancy in the plane or solution to keep it simple and concise
- Provide 1-2 specific suggestions for how the solution plan could be improved.
- propose one step plan, optional
- propose next sub plan along with answer, to a plan optional
- propose how to answer a step of plan again, optional
- propose rephrase question or subquestion to a plan, optional


** Multi-dimensional Scoring **
Dimensions to score:
Reasoning: <How likely is it that the conclusion will be overturned?>
Score Given: <score ,[-30-0]> 

Solution Accuracy and Error-Free:
Reasoning: <How accurate and error-free is the solution?>
Score Given: <score ,[0-10]> 

Solution Correctness:
Reasoning: <How thoroughly does it address all aspects of the problem?>
Score Given: <score ,[0-10]> 

Solution Completeness:
Reasoning: <how complete is the solution?>
Score Given: <score ,[0-10]> 

Solution Clarity: 
Reasoning: <How clear and easy to understand is the explanation?>
Score Given: <score ,[0-10]> 

Solution Efficiency: 
Reasoning: <How optimal is the approach in terms of time/resource usage?>
Score Given: <score ,[0-10]> 

Solution Redundancy: 
Reasoning: < Is there one or more step can be  deleted or simplified>
Score Given: <score ,[-20-0]> 

** Overall Evaluation **
Sum Score Calculation: 
- calculate step by step to get the sum of the above scores 
Overall Evaluation: <display calculated score > 
`)

	var bestAnswerNode *TreeNode = MCTSTrajectory.BestScoreNode()
	for i := 0; i < TryNum; i++ {
		NewBestTrail := bestAnswerNode.NewChild()
		if NewBestTrail.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, bestAnswerNode.Solution, bestAnswerNode.Refinement("improvement Suggestions"), SolverMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if NewBestTrail.Verification, err = models.SLM2.AskLLM(0.7, false, MCTSTrajectory.UserMsg, NewBestTrail.Solution, VerifierMessege); err != nil {
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
