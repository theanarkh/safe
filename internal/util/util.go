package util

import "fmt"

func SafeCall(f func()) {
	if f == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic happen in handler: %v\n", err)
		}
	}()
	f()
}
