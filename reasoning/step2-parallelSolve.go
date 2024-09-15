package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

func (node *TreeNode) ParalleSolve(difficulty float64) (err error) {

	//generateTrajectory

	// Based on the difficulty, recommend a compute strategy:
	// - For easier problems (1-2): Evaluate the strengths and weaknesses of solutions
	// - For moderate problems (3): Evaluate the strengths and weaknesses of solutions. Propose alternative approaches or solutions
	// - For harder problems (4-5): Propose alternative approaches or solutions. Evaluate the strengths and weaknesses of each",

	// ## Solution Refinement:
	// <Refinement: Evaluate the strengths and weaknesses of solutions>
	// ## Alternative Solution:
	// <Alternative Solution here, when difficulty >= 3>
	// ## Alternative Solution Refinement:
	// <Alternative Solution here, when difficulty is 5>
	numParalellUnfolding := int(2 * difficulty)
	for i := len(MCTSTrajectory.Children()); i < numParalellUnfolding; i++ {
		childNode := MCTSTrajectory.NewChild()
		childNode.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, models.UserMsg(`Now, follow these steps:
1. Problem Analysis:
Analyze the given problem. Identify key components, constraints, and potential approaches.

2. Plan Proposal:
- Devise a step-by-step plan to solve the problem. (don't actually start solving yet, just make a plan)

3. Solution to the question:
- Use Chain of Thought reasoning to work through the plan and write the full solution within thinking.
## Solution: <...>

4. Solution Refinement:
Evaluate the strengths and weaknesses of Plan Proposal
`))
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println("The Solution is: ", childNode.Solution.Content)
		CheckerMesseges := models.UserMsg(`Given a question and Problem Analysis, Plan Proposal, Solution and it's refinement proposal, Follow these steps:
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
Overall Evaluation: <display calculated Sum Score > 
`)
		if childNode.Verification, err = models.SLM2.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, childNode.Solution, CheckerMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		childNode.Save()
	}
	return nil
}
