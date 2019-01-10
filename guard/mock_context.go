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

// Ctx return origin context from the mock context, not from the parent context.
// Otherwise, Control from the mock context
// will be canceled after the parent context.
func (ctx MockContext) Ctx() context.Context {
	return ctx.Mock.Ctx()
}

// Override UpdateCtx for prevent replacing the parent context.
// Otherwise, it will cancel the parent context when the mock context cancels.
func (ctx MockContext) UpdateCtx(newCtx context.Context) {
	ctx.Mock.UpdateCtx(newCtx)
}

// Done returns the channel of the Mock context.
func (ctx MockContext) Done() <-chan struct{} {
	return ctx.Mock.Done()
}
