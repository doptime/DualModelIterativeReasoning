package main

import (
	"DualModelIterativeReasoning/reasoning"
)

func main() {
	mp, _ := reasoning.KeyTreeNode.HGetAll()

	for _, v := range mp {
		reasoning.NodesMap.Set(v.Id, v)
	}
	if node, ok := reasoning.NodesMap.Get("root"); ok {
		reasoning.MCTSTrajectory = node
	}
	if reasoning.NodesMap.Count() == 0 {
		reasoning.NodesMap.Set("root", reasoning.MCTSTrajectory)
	}
	// Perform reasoning
	reasoning.MCTSTrajectory.DualModelIterativeResoning()

}
