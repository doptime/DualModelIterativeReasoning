package reasoning

import (
	"DualModelIterativeReasoning/models"
	"fmt"
)

var MCTSTrajectory = &TreeNode{
	Id:      "root",
	UserMsg: models.MsgOfUser("Discussion: what's the top 15 things (foods, drugs, excercises, ...) should be in the list, In order to maximize life expectancy. while a man have family diabetes background but not have diabetes currently ?"),
	//UserMsg: models.MsgOfUser("Discussion:Here's several food additives in the list: taurine, beta-glucan, quercetin, turmeric-black pepper tablets, zinc gluconate tablets, brewer's yeast tablets, VD oil, fish oil, vitamin k2, Broad spectrum probiotics, Berberine,Magnesium. In order to maximize life expectancy, and contrains ingredients no more than 15 , after remove and adds to the list, what is the final ingredients to be in the list?"),
	//UserMsg: models.UserMsg("If a layer of material of uniform thickness is applied on the surface of an ellipsoid, will the surface still be an strict mathematically ellipsoid after coating?"),
}
var SysPromptBasic = models.MsgOfUser("You are a world-class powerfull AI system, cooperative, innovative, reflective and helpfull, capable of complex reasoning. Together with another AI model, you are solving problems through structured collaboration. ")

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
	err = node.SequentialSolve(30)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	var bestAnswerNode *TreeNode = node.BestScoreNode()
	fmt.Println("The Best Asnwer is: ", bestAnswerNode.Solution)
	return nil
}
