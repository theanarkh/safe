package errgroup

import (
	"context"
	"fmt"

	egroup "golang.org/x/sync/errgroup"
)

type handlerFunc func(err any) error

type Group struct {
	egroup.Group
	cancel  func(error)
	handler handlerFunc
}

type Option func(*Group)

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

func WithHandler(handler handlerFunc) func(e *Group) {
	return func(e *Group) {
		e.handler = handler
	}
}

func New(options ...Option) *Group {
	eg := &Group{}
	for _, option := range options {
		option(eg)
	}
	return eg
}

func (g *Group) handlePanic(e *error) {
	if err := recover(); err != nil {
		*e = fmt.Errorf("panic happen: %v", err)
		if g.handler == nil {
			fmt.Println(e)
		} else {
			result := g.handler(err)
			if result != nil {
				*e = result
			}
		}
	}
}

func (g *Group) Go(f func() error) {
	g.Group.Go(func() (err error) {
		defer g.handlePanic(&err)
		return f()
	})
}

func (g *Group) TryGo(f func() error) bool {
	return g.Group.TryGo(func() (err error) {
		defer g.handlePanic(&err)
		return f()
	})
}

func (g *Group) Wait() (err error) {
	defer g.handlePanic(&err)
	err = g.Group.Wait()
	if g.cancel != nil {
		g.cancel(err)
	}
	return err
}
