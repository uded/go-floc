package guard

import (
	"github.com/uded/go-floc"
)

/*
Cancel cancels execution of the flow with the data given.
*/
func Cancel(data interface{}) floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		ctrl.Cancel(data)
		return nil
	}
}
