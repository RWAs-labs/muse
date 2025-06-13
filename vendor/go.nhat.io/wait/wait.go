// Package wait provides functionalities for waiting for a condition.
package wait

import (
	"context"
	"time"
)

// Waiter is an interface that waits for a condition.
//
//go:generate mockery --name Waiter --output ./mock --outpkg wait --filename waiter.go
type Waiter interface {
	Wait(ctx context.Context) error
}

// NoWait is a Waiter that does not wait.
var NoWait = noWait{}

var _ Waiter = (*noWait)(nil)

// noWait is a Waiter.noWait.
type noWait struct{}

// Wait does not wait.
func (noWait) Wait(context.Context) error {
	return nil
}

var _ Waiter = (*Func)(nil)

// Func is a function that delays the execution.
type Func func(ctx context.Context) error

// Wait waits for the function to return.
func (f Func) Wait(ctx context.Context) error {
	return f(ctx)
}

var _ Waiter = (*ForDuration)(nil)

// ForDuration is a Waiter that waits for a duration.
type ForDuration time.Duration

// Wait waits for the duration unless the context is canceled.
func (d ForDuration) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(time.Duration(d)):
		return nil
	}
}

var _ Waiter = (*Until)(nil)

// Until is a Waiter that waits until a time.
type Until time.Time

// Wait waits until the time unless the context is canceled.
func (t Until) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(time.Until(time.Time(t))):
		return nil
	}
}

var _ Waiter = (*ForSignal)(nil)

// ForSignal is a Waiter that waits for a signal.
type ForSignal <-chan time.Time

// Wait waits for a signal unless the context is canceled.
func (s ForSignal) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-s:
		return nil
	}
}
