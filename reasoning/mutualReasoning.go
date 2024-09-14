package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"

	"github.com/doptime/doptime/db"
)

var MCTSTrajectory = &TreeNode{
	Id: "root",
	// UserMsg: models.UserMsg("Discussion:Here's several food additives: taurine, beta-glucan, quercetin, turmeric-black pepper tablets, tea polyphenols, selenium yeast tablets, resveratrol, zinc gluconate tablets, VC tablets, brewer's yeast tablets. In order to maximize life expectancy, which one ingredients should  delete and which new ingredient to add?"),
	UserMsg: models.UserMsg("If a layer of material of uniform thickness is applied on the surface of an ellipsoid, will the surface still be an ellipsoid after coating?"),
}
var SysPromptBasic = models.SystemMsg("You are a world-class powerfull AI system, cooperative, innovative, reflective and helpfull, capable of complex reasoning. Together with another AI model, you are solving problems through structured collaboration. ")

func DualModelIterativeResoning() (err error) {

	//generateTrajectory
	SoluteRequestFirstTime := models.UserMsg(`Follow these steps:

1. Problem Analysis:
Analyze the given problem. Identify key components, constraints, and potential approaches.

2. Initial Proposal:
Provide an initial solution or approach to the problem. Be detailed and explain your reasoning.

3. Difficulty Estimation:
Estimate the difficulty of this problem on a scale of 1-5, integer, where:
1 = Very Easy
2 = Easy
3 = Moderate
4 = Difficult
5 = Very Difficult
;using format:
## difficulty: ...

Explain your reasoning for this estimation.

4. give solution to the question;
## Solution: ...

5. Solution Refinement:
Based on the difficulty, recommend a compute strategy:
- For easier problems (1-2): Evaluate the strengths and weaknesses of solutions
- For moderate problems (3): Evaluate the strengths and weaknesses of solutions. Propose alternative approaches or solutions
- For harder problems (4-5): Propose alternative approaches or solutions. Evaluate the strengths and weaknesses of each",

## Solution Refinement: 
<Refinement: Evaluate the strengths and weaknesses of solutions>
## Alternative Solution:
<Alternative Solution here, when difficulty >= 3>
## Alternative Solution Refinement:
<Alternative Solution here, when difficulty is 5>

`)
	for i, L := len(MCTSTrajectory.Children()), 4; i < L; i++ {
		childNode := MCTSTrajectory.NewChild(db.NanoId(8))
		childNode.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, SoluteRequestFirstTime)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println("The Solution is: ", childNode.Solution.Content)
		CheckerMesseges := models.UserMsg(`Given a question and a answer, Follow these steps:
1. Solution Refinement:
- <item: identify potential improvements or flaws in the current solution>
- <item: propose specific changes or additions>


2. Solution Sore:
- overall score: <overall score for the current solution or proposed alternatives, real value, between [-1,1] >
-1 = Incorrect or significantly flawed
0 = Partially correct but needs improvement
1 = Correct and well-reasoned
- scoring reasoning.
`)
		if childNode.Verification, err = models.SLM2.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, childNode.Solution, CheckerMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		childNode.Save()
	}

	SolverMesseges := models.UserMsg(`Given a question, and solution together with refinement ideas. Follow these steps tofurther solve the problem:
1. Problem Analysis:
Analyze the given problem, solution and refinement ideas. Identify key components, constraints, and potential approaches.


2. Solution Proposal:
<propose solution here based on the former version and refinement ideas>
`)
	VerifierMessege := models.UserMsg(`Given a question, and solution (may be together with refinement ideas). Follow these steps tofurther solve the problem:

1. Solution Refinement:
- <item: identify potential important correctness in the current solution>
- <item: identify potential improvements in the current solution>
- <item: identify flaws in the current solution>
- <item: propose specific changes or additions>

2. Solution Scoring:
- overall score: <overall score for the current solution or proposed alternatives, real value, between [-1,1] >
-1 = Incorrect or significantly flawed
0 = Partially correct but needs improvement
1 = Correct and well-reasoned
- scoring reasoning.
`)

	for i := 0; i < 4; i++ {
		var bestAnswerNode *TreeNode = MCTSTrajectory.BestScoreNode()
		NewBestTrail := bestAnswerNode.NewChild(db.NanoId(8))
		if NewBestTrail.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, bestAnswerNode.Solution, bestAnswerNode.Refinement("Refinement:"), SolverMesseges); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if NewBestTrail.Verification, err = models.SLM2.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, NewBestTrail.Solution, VerifierMessege); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		NewBestTrail.Save()
	}
	var bestAnswerNode *TreeNode = MCTSTrajectory.BestScoreNode()
	fmt.Println("The Best Asnwer is: ", bestAnswerNode.Solution)
	return nil
}
