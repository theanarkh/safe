package errgroup

import (
	"errors"
	"sync/atomic"
	"testing"
)

func TestGroupNew(t *testing.T) {
	g := New()
	g.Go(func() error {
		return nil
	})
	err := g.Wait()
	if err != nil {
		t.Fatalf("should return nil")
	}
}

func TestGroupNewPanic(t *testing.T) {
	var count atomic.Int32
	g := New(
		WithHandler(func(err error) {
			if err == nil {
				t.Fatalf("should be error")
			}
			count.Add(1)
		}),
	)
	g.Go(func() error {
		panic("panic test1")
	})
	g.Go(func() error {
		panic("panic test2")
	})
	err := g.Wait()
	if err == nil {
		t.Fatalf("should return error")
	}
	if count.Load() != 2 {
		t.Fatalf("should call handler twice")
	}
}
func TestGroupNewError(t *testing.T) {
	var count atomic.Int32
	g := New(
		WithHandler(func(err error) {
			if err == nil {
				t.Fatalf("should be error")
			}
			count.Add(1)
		}),
	)
	g.Go(func() error {
		return errors.New("error test1")
	})
	g.Go(func() error {
		return errors.New("error test1")
	})
	err := g.Wait()
	if err == nil {
		t.Fatalf("should return error")
	}
	if count.Load() != 2 {
		t.Fatalf("should call handler twice")
	}
}
