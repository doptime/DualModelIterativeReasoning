package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

var MCTSTrajectory = &TreeNode{
	Id: "root",
	// UserMsg: models.UserMsg("Discussion:Here's several food additives: taurine, beta-glucan, quercetin, turmeric-black pepper tablets, tea polyphenols, selenium yeast tablets, resveratrol, zinc gluconate tablets, VC tablets, brewer's yeast tablets. In order to maximize life expectancy, which one ingredients should  delete and which new ingredient to add?"),
	UserMsg: models.UserMsg("If a layer of material of uniform thickness is applied on the surface of an ellipsoid, will the surface still be an ellipsoid after coating?"),
}
var SysPromptBasic = models.SystemMsg("You are a world-class powerfull AI system, cooperative, innovative, reflective and helpfull, capable of complex reasoning. Together with another AI model, you are solving problems through structured collaboration. ")

func (node *TreeNode) DualModelIterativeResoning() (err error) {
	var difficulty float64

	//先搭建完善求解架构。然后再求解。
	difficulty, err = node.GetDifficulty()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	err = node.ParalleSolve(difficulty)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	err = node.SequentialSolve(10)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	var bestAnswerNode *TreeNode = node.BestScoreNode()
	fmt.Println("The Best Asnwer is: ", bestAnswerNode.Solution)
	return nil
}
