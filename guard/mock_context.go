package guard

import (
	"context"

	"gopkg.in/devishot/go-floc.v2"
)

// Mock context which propagates all calls to the parent context
// but Done() returns Mock channel.
type MockContext struct {
	floc.Context
	Mock floc.Context
}

// Release releases the Mock context.
func (ctx MockContext) Release() {
	ctx.Mock.Release()
}

// If we prevent replacing the parent context, then it won't be canceled
func (ctx MockContext) UpdateCtx(newCtx context.Context) {
	ctx.Mock.UpdateCtx(newCtx)
}

// Done returns the channel of the Mock context.
func (ctx MockContext) Done() <-chan struct{} {
	return ctx.Mock.Done()
}
