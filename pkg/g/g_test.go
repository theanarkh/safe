package g

import (
	"sync"
	"testing"
)

func TestGo(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	Go(func() {
		defer wg.Done()
	})
	wg.Wait()
}

func TestGoRecover(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	handler1 := func(err error) {
		if err == nil {
			t.Fatalf("should propagate panic")
		}
	}
	handler2 := func(err error) {
		wg.Done()
	}
	Go(func() { panic("panic test") }, handler1, handler2)
	wg.Wait()
}

func TestG(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	var g G
	g.Go(func() {
		defer wg.Done()
	})
	wg.Wait()
}

func TestGRecover(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	var g1 G
	g1.Go(func() {
		defer wg.Done()
		panic("panic test")
	})
	wg.Wait()
}

func TestGNew(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	g := New()
	g.Go(func() {
		defer wg.Done()
	})
	wg.Wait()
}

func TestGNewRecover(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	g := New(
		WithHandler(func(err error) {
			if err == nil {
				t.Fatalf("should propagate panic")
			}
		}),
		WithHandler(func(err error) {
			wg.Done()
		}),
	)
	g.Go(func() { panic("panic test") })
	wg.Wait()
}

func TestGoHandlerPanic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	handler := func(err error) {
		defer wg.Done()
		panic("handler panic")
	}
	Go(func() { panic("panic test") }, handler)
	wg.Wait()
}
