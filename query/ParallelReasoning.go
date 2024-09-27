package query

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func AskLLMParallelly(query ...*Query) (err error) {
	// exeInGroup := func(g *errgroup.Group, node *TreeNode) {
	// 	g.Go(func() (err error) {
	// 		return node.Solute()
	// 	})
	// }
	g, _ := errgroup.WithContext(context.Background())
	for _, q := range query {
		//exeInGroup(g, q)

		g.Go(func() (err error) {
			return q.Solute()
		})
	}
	err = g.Wait()
	return err
}
