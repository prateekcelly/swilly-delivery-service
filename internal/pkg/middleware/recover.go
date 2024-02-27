package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
)

type RecoverableFunc func(context.Context) error

func ProcessWithRecovery(ctx context.Context, fn RecoverableFunc) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			debug.PrintStack()
		}
	}()

	return fn(ctx)
}
