package guard

import "gopkg.in/devishot/go-floc.v2"

type DuplicateContext struct {
  floc.Context
  copy floc.Context
}

func NewDuplicateContext(ctx floc.Context) floc.Context {
  newCtx := floc.BorrowContext(ctx.Ctx())
  return &DuplicateContext{Context: newCtx, copy: ctx}
}

func (c *DuplicateContext) AddValue(key, value interface{}) {
  c.Context.AddValue(key, value)
  c.copy.AddValue(key, value)
}
