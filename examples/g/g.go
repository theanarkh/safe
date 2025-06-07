package main

import (
	"fmt"
	"sync"

	"github.com/theanarkh/safe/pkg/g"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)
	{
		g.Go(func() {
			defer wg.Done()
		})
	}

	{
		var g g.G
		g.Go(func() {
			defer wg.Done()
			panic("panic test")
		})
	}

	{
		g := g.New(
			g.WithHandler(func(err error) {
				defer wg.Done()
			}),
		)
		g.Go(func() {
			panic("panic test1")
		})
		g.Go(func() {
			panic("panic test2")
		})
	}
	wg.Wait()
	fmt.Println("done")
}
