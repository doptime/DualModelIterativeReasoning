package reasoning

// import (
// 	"DualModelIterativeReasoning/message"
// 	"DualModelIterativeReasoning/models"
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/samber/lo"
// 	"golang.org/x/sync/errgroup"
// )

// func (node *TreeNode) ParalleSolve(difficulty float64) (err error) {
// 	numParallelUnfolding := int(2 * difficulty)
// 	genSolution := message.User(`Now, follow these steps:
// 1. Problem Analysis:
// Analyze the given problem. Identify key components, constraints, and potential approaches.

// 2. Plan Proposal:
// - Devise a step-by-step plan to solve the problem. (don't actually start solving yet, just make a plan)

// 3. Solution to the question:
// - Use Chain of Thought reasoning to work through the plan and write the full solution within thinking.
// ## Solution: <...>

// 4. Solution Refinement:
// Evaluate the strengths and weaknesses of Plan Proposal
// 	`)
// 	g, _ := errgroup.WithContext(context.Background())
// 	for i := len(node.Children()); i < numParallelUnfolding; i++ {
// 		g.Go(func() (err error) {
// 			childNode := node.NewChild()
// 			childNode.BestSolution = true
// 			childNode.Solution, err = models.SLM1.AskLLM(0.7, false, SysPromptBasic, node.UserMsg, genSolution)
// 			if err == nil {
// 				fmt.Println("The Solution is: ", childNode.Solution.Content)
// 				childNode.Save()
// 			}
// 			return err
// 		})
// 	}
// 	err = g.Wait()

// 	children := node.Children()
// 	node.ParalleSolveChooseBestAnswer(children)
// 	return err
// }

// func (rootNode *TreeNode) ParalleSolveChooseBestAnswer(Nodes []*TreeNode) (bestAnswerNode *TreeNode) {
// 	for true {
// 		var leftNodes = lo.Filter(Nodes, func(node *TreeNode, ind int) bool {
// 			return node.BestSolution
// 		})
// 		if len(leftNodes) < 1 {
// 			break
// 		}
// 		for j := 0; j < len(leftNodes); j += 2 {
// 			solution1, solution2 := leftNodes[j], leftNodes[j+1]
// 			g, _ := errgroup.WithContext(context.Background())
// 			g.Go(func() (err error) {
// 				Solution, err := models.SLM2.AskLLM(0.7, false, message.User(SysPromptBasic.Content+"\r"+rootNode.UserMsg.Content), message.Assistant("ok. this is the original question."),
// 					message.User("this is Solution1 to the question:"+solution1.Solution.Content), message.Assistant("ok. this is the solution1 to the question."),
// 					message.User("this is Solution2 to the question:"+solution2.Solution.Content), message.Assistant("ok. this is the solution2 to the question."),
// 					message.User(`Now, follow these steps to choose the best solution:
// 1. Analyze solution1, reasoning the merits and weakness of solution1
// 2. Analyze solution2, reasoning the merits and weakness of solution2
// 3. Choose One best Solution to kept as Final Solution:
// - reasoning which One Solution to choose as Final Solution

// 4. The best solution choosen is: <Solution1 or <Solution2>> `))
// 				if err != nil {
// 					return err
// 				}
// 				solution2isBest := strings.Index(Solution.Content, "olution2") > strings.Index(Solution.Content, "olution1")
// 				if solution2isBest {
// 					solution1.BestSolution = false
// 					solution1.Save()
// 				} else {
// 					solution2.BestSolution = false
// 					solution2.Save()
// 				}

// 				return err
// 			})
// 		}

// 	}
// }
