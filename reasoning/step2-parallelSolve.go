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
		childNode.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, models.MsgOfUser(`Now, follow these steps:
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
		if childNode.Verification, err = models.SLM2.AskLLM(0.7, false, models.MsgOfUser(SysPromptBasic.Content+"\r"+MCTSTrajectory.UserMsg.Content), childNode.Solution, models.MsgOfUser(`Given a question, Problem Analysis, Plan Proposal, Solution, and its refinement proposal, follow these steps:
		
		** Critical Analysis and Improvement Suggestions **
		1. Alternative Perspective:
		   - Propose the strongest possible argument or approach that could overturn the current conclusion.
		   - Explain the reasoning behind this alternative perspective.
		
		2. Redundancy and Simplification:
		   - Identify any redundant or unnecessary steps in the plan or solution.
		   - Suggest how to simplify the approach while maintaining its effectiveness.
		
		3. Specific Improvements:
		   - Provide 2-3 specific, actionable suggestions to improve the solution plan.
		   - For each suggestion, explain its potential impact on the overall solution.
		
		4. Optional Refinements (choose 1-2 if applicable):
		   - Propose a single-step plan to address a weakness in the current solution.
		   - Suggest a sub-plan along with its potential answer to deepen the solution.
		   - Recommend how to rephrase a step or sub-question for clarity or better focus.
		
		** Multi-dimensional Scoring **
		Score each dimension from 0-10, where 0 is the worst and 10 is the best. Provide a brief justification for each score.
		
		1. Reasoning Quality:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		2. Solution Accuracy:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		3. Problem Coverage:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		4. Clarity of Explanation:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		5. Efficiency of Approach:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		6. Self-Criticism and Alternative Perspectives:
		   Justification: <1-2 sentences explaining the score>
		   Score: [0-10]
		
		** Overall Evaluation **
		Total Score Calculation:
		- Sum up the scores from all six dimensions
		- Maximum possible score: 60
		
		Overall Score: <display calculated total score>
		Percentage: <(Total Score / 60) * 100%>
		
		Interpretation:
		- 90-100%: Excellent solution with strong self-criticism
		- 70-89%: Good solution with some self-reflection, minor improvements needed
		- 50-69%: Average solution, more critical thinking and alternative perspectives needed
		- Below 50%: Poor solution, lacks self-criticism and alternative viewpoints
		
		Final Recommendations:
		- Identify the lowest-scoring dimension and suggest a focused approach to improve it in the next iteration.
		- Provide a concise summary of the most crucial improvements needed, emphasizing self-criticism and consideration of alternative perspectives.
		`)); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		childNode.Save()
	}
	return nil
}
