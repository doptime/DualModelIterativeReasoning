package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

func (node *TreeNode) GetDifficulty() (difficulty float64, err error) {

	DifficultyEstimation := models.UserMsg(`Given the question. Follow these steps:
## Difficulty explanations:
Explain your reasoning for this estimation on difficulty.

## Difficulty Estimation:
difficulty level: <difficulty value>
Estimate the difficulty of this problem on a scale of 1-5, integer, where:
1 = Very Easy
2 = Easy
3 = Moderate
4 = Difficult
5 = Very Difficult
`)
	if MCTSTrajectory.Solution == nil {
		MCTSTrajectory.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, MCTSTrajectory.UserMsg, DifficultyEstimation)
		if err != nil {
			fmt.Println("Error: ", err)
			return 0, err
		} else {
			MCTSTrajectory.Save()
		}
	}

	//先搭建完善求解架构。然后再求解。
	MCTSTrajectory.Difficulty, err = ReadFloatAfterTag("level:", MCTSTrajectory.Solution.Content)
	return MCTSTrajectory.Difficulty, err
}
