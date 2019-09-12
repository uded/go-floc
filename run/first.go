package run

import (
	"github.com/uded/go-floc"
	"github.com/uded/go-floc/guard"
)

/*
First runs jobs in their own goroutines and waits until first of them finish.

Summary:
	- Run jobs in goroutines : YES
	- Wait all jobs finish   : NO
	- Run order              : PARALLEL

Diagram:
    +-->[JOB_1]--+
    |            |
  --+-->  ..     +-->
    |            |
    +-->[JOB_N]--+
*/
func First(jobs ...floc.Job) floc.Job {
	return func(ctx floc.Context, ctrl floc.Control) error {
		// Do not start parallel jobs if the execution is finished
		if ctrl.IsFinished() {
			return nil
		}

		newCtx := guard.NewDuplicateContext(ctx)

		mockCtx := guard.MockContext{
			Context: newCtx,
			Mock:    floc.NewContext(),
		}
		defer mockCtx.Release()

		mockCtrl := floc.NewControl(mockCtx)
		defer mockCtrl.Release()

		// Run jobs in parallel
		for _, job := range jobs {
			// Run the job in it's own goroutine
			go func(job floc.Job) {
				var err error
				err = job(mockCtx, mockCtrl)
				handleResult(mockCtrl, err)
			}(job)
		}

		// Wait until first jobs done
		<-newCtx.Done()

		res, data, err := mockCtrl.Result()
		switch res {
		case floc.Completed:
			mockCtrl.Cancel(nil)
			// Continue current flow
		case floc.Canceled:
			ctrl.Cancel(data)
		case floc.Failed:
			ctrl.Fail(data, err)
		}

		return nil
	}
}
