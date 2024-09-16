package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

func (node *TreeNode) GetDifficulty() (difficulty float64, err error) {

	DifficultyEstimation := models.MsgOfUser(`Given the question. Follow these steps:
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
	if MCTSTrajectory.Solution != nil && len(MCTSTrajectory.Solution.Content) > 0 {
		MCTSTrajectory.Difficulty, _ = ReadFloatAfterTag(MCTSTrajectory.Solution.Content, "level:", "level of")
	}
	if MCTSTrajectory.Difficulty == 0 {
		MCTSTrajectory.Solution, err = models.SLM1.AskLLM(0.7, false, models.MsgOfUser(SysPromptBasic.Content+";\n"+MCTSTrajectory.UserMsg.Content+";\n"+DifficultyEstimation.Content))
		if err != nil {
			fmt.Println("Error: ", err)
			return 0, err
		} else {
			//先搭建完善求解架构。然后再求解。
			MCTSTrajectory.Difficulty, err = ReadFloatAfterTag(MCTSTrajectory.Solution.Content, "level:", "level of")
			MCTSTrajectory.Save()
		}
	}

	return MCTSTrajectory.Difficulty, err
}
