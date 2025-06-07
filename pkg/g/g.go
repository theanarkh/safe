package g

import (
	"fmt"
	"runtime/debug"

	"github.com/theanarkh/safe/internal/util"
)

type Handler func(error)

type Option func(g *G)

type G struct {
	handler Handler
}

func Go(f func(), handlers ...Handler) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if len(handlers) != 0 {
					e := fmt.Errorf("panic happen: %v, call stack: %s", err, string(debug.Stack()))
					for _, handler := range handlers {
						if handler != nil {
							util.SafeCall(func() {
								handler(e)
							})
						}
					}
				}
			}
		}()
		f()
	}()
}

func WithHandler(h func(err error)) Option {
	return func(g *G) {
		g.handler = h
	}
}

func New(opts ...Option) *G {
	g := &G{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *G) Go(f func()) {
	Go(f, g.handler)
}
