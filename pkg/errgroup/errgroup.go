package errgroup

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/theanarkh/safe/internal/util"
	egroup "golang.org/x/sync/errgroup"
)

type Handler func(err error)
type Option func(*Group)

type Group struct {
	egroup.Group
	cancel  func(error)
	handler Handler
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

func WithHandler(handler Handler) func(e *Group) {
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

func (g *Group) callHandler(err error) {
	if g.handler != nil {
		util.SafeCall(func() {
			g.handler(err)
		})
	}
}

func (g *Group) recover(e *error) {
	if err := recover(); err != nil {
		*e = fmt.Errorf("panic happen: %v, call stack: %s", err, string(debug.Stack()))
		g.callHandler(*e)
	}
}

func (g *Group) Go(f func() error) {
	g.Group.Go(func() (err error) {
		defer g.recover(&err)
		err = f()
		if err != nil {
			g.callHandler(err)
		}
		return
	})
}

func (g *Group) TryGo(f func() error) bool {
	return g.Group.TryGo(func() (err error) {
		defer g.recover(&err)
		err = f()
		if err != nil {
			g.callHandler(err)
		}
		return
	})
}

func (g *Group) Wait() (err error) {
	defer g.recover(&err)
	err = g.Group.Wait()
	if g.cancel != nil {
		g.cancel(err)
	}
	return
}
