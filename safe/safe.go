package safe

import "fmt"

func Go(f func(), handler func(err any)) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if handler == nil {
					fmt.Println(fmt.Sprintf("panic happen: %v", err))
				} else {
					handler(err)
				}
			}
		}()
		f()
	}()
}
