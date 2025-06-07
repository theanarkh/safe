package main

import (
	"fmt"

	"github.com/theanarkh/safe/pkg/errgroup"
)

func main() {
	{
		var eg errgroup.Group
		eg.Go(func() error {
			return nil
		})
		err := eg.Wait()
		fmt.Println(err)
	}
	{
		var eg errgroup.Group
		eg.Go(func() error {
			panic("panic test1")
		})
		err := eg.Wait()
		fmt.Println(err)
	}
	{
		eg := errgroup.New(
			errgroup.WithHandler(func(err error) {
				fmt.Printf("panic:%s\n", err.Error())
			}),
		)
		eg.Go(func() error {
			panic("panic test2")
		})
		eg.Go(func() error {
			panic("panic test3")
		})
		err := eg.Wait()
		fmt.Println("final error: %s" + err.Error())
	}
}
