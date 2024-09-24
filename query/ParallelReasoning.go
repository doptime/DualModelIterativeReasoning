package query

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func ParallelReasnoing(querys ...[]*TreeNode) (err error) {
	g, _ := errgroup.WithContext(context.Background())
	for i, query := range querys {
		g.Go(func() (err error) {
			return query[i].Solute()
		})
	}
	err = g.Wait()
	return err
}
