package main

import (
	"DualModelIterativeReasoning/message"
	"DualModelIterativeReasoning/reasoning"
)

// Perform reasoning
var MCTSTrajectory = &reasoning.TreeNode{
	Id:      "root",
	UserMsg: message.User("Write a story about the plot tension and the plot. Given clues. Midi, both hands, stolen goods, subnets, random eating, small yellow car. The plot iteration is ups and downs, which is very attractive."),
	// UserMsg: message.User("Discussion: what's the top 15 things (foods, drugs, excercises, ...) should be in the list, In order to maximize life expectancy. while a man have family diabetes background but not have diabetes currently ?"),
	// UserMsg: message.User("Discussion:Here's several food additives in the list: taurine, beta-glucan, quercetin, turmeric-black pepper tablets, zinc gluconate tablets, brewer's yeast tablets, VD oil, fish oil, vitamin k2, Broad spectrum probiotics, Berberine,Magnesium. In order to maximize life expectancy, and contrains ingredients no more than 15 , after remove and adds to the list, what is the final ingredients to be in the list?"),
	// UserMsg: models.UserMsg("If a layer of material of uniform thickness is applied on the surface of an ellipsoid, will the surface still be an strict mathematically ellipsoid after coating?"),
}

func main() {
	mp, _ := reasoning.KeyTreeNode.HGetAll()

	for _, v := range mp {
		reasoning.NodesMap.Set(v.Id, v)
	}
	if node, ok := reasoning.NodesMap.Get("root"); ok {
		MCTSTrajectory = node
	}
	if reasoning.NodesMap.Count() == 0 {
		reasoning.NodesMap.Set("root", MCTSTrajectory)
	}

	//step1 get difficulty
	MCTSTrajectory.ParalleBeamSearchUsingDualModelIterativeReasoning(20)

}
