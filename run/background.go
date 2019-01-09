package run

import (
	"gopkg.in/devishot/go-floc.v2"
	"gopkg.in/devishot/go-floc.v2/guard"
)

/*
Background starts the job in it's own goroutine. The function does not
track the lifecycle of the job started and does no synchronization with it
therefore the job running in background may remain active even if the flow
is finished. The function assumes the job is aware of synchronization and termination
of it is implemented outside. The job can not change the flow state.

	floc.Run(run.Background(
		func(ctx floc.Context, ctrl floc.Control) error {
			for !ctrl.IsFinished() {
				fmt.Println(time.Now())
			}

			return nil
		}
	})

Summary:
	- Run jobs in goroutines : YES
	- Wait all jobs finish   : NO
	- Run order              : SINGLE

Diagram:
  --+----------->
    |
    +-->[JOB]
*/
func Background(job floc.Job) floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		// Do not start the job if the flow is finished
		if ctrl.IsFinished() {
			return nil
		}

		// Create new context
		mockCtx := guard.MockContext{
			Context: ctx,
			Mock:    floc.NewContext(),
		}
		mockCtrl := floc.NewControl(mockCtx)

		// Run the job in background
		go func(job floc.Job) {
			defer mockCtx.Release()
			defer mockCtrl.Release()

			err := job(mockCtx, mockCtrl)
			handleResult(mockCtrl, err)
		}(job)

		return nil
	}
}
