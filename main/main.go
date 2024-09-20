package main

import (
	"DualModelIterativeReasoning/message"
	"DualModelIterativeReasoning/reasoning"
)

// Perform reasoning
var MCTSTrajectory = &reasoning.TreeNode{
	Id: "root",
	UserMsg: message.User(`I have a 9 yrs old daughter, I want's help here with her using a funny | interesting | breath taking | deep-diving | emotion arousing story. 
Remember, The Most Important thing is building the experience. If can not, Others Fade Away, because she's somehow formidable with her work.

This is the math problem:
Duckling 1000 grams, lamb 60 kilograms, puppy 8 kilograms, calf 130 kilograms, shark 2 tons
1. Among the above animals, () is the heaviest and () is the lightest
2. The puppy is () grams heavier than the duckling, and the lamb is () kilograms lighter than the calf
3. A calf and a lamb weigh () kilograms in total, and () sharks like this weigh 10 tons in total.

`),
}

func main() {
	mp, _ := reasoning.KeyTreeNode.HGetAll()

	for _, v := range mp {
		reasoning.NodesMap.Set(v.Id, v)
	}
	reasoning.NodesMap.Set("root", MCTSTrajectory)
	if node, ok := reasoning.NodesMap.Get("root"); ok {
		MCTSTrajectory = node
	}
	if reasoning.NodesMap.Count() == 0 {
		reasoning.NodesMap.Set("root", MCTSTrajectory)
	}

	//step1 get difficulty
	MCTSTrajectory.ParalleBeamSearchUsingDualModelIterativeReasoning(20)

}
