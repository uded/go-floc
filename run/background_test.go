package run

import (
	"sync"
	"testing"
	"time"

	"fmt"

	"gopkg.in/devishot/go-floc.v2"
)

func TestBackground_AlreadyFinished(t *testing.T) {
	ctx := floc.NewContext()
	defer ctx.Release()

	ctrl := floc.NewControl(ctx)
	defer ctrl.Release()

	flow := Background(complete(nil))

	ctrl.Cancel(nil)

	result, _, _ := floc.RunWith(ctx, ctrl, flow)
	if !result.IsCanceled() {
		t.Fatalf("%s expects result to be Canceled but has %s", t.Name(), result.String())
	}
}

func TestBackground_ParentFlowFinished(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	flow := Sequence(
		Background(
			Delay(time.Millisecond,
				Sequence(
					func(ctx floc.Context, ctrl floc.Control) error {
						wg.Done()
						return nil
					},
					complete(nil),
				),
			),
		),
		complete(nil),
	)

	result, _, _ := floc.Run(flow)
	if !result.IsFinished() {
		t.Fatalf("%s expects result to be Finished but has %s", t.Name(), result.String())
	}

	wg.Wait()
}

func TestBackground_Completed(t *testing.T) {
	flow := Sequence(
		Background(cancel(nil)),
		Delay(time.Millisecond, complete(nil)),
	)
	result, data, err := floc.Run(flow)
	if !result.IsCompleted() {
		t.Fatalf("%s expects result to be Completed but has %s", t.Name(), result.String())
	} else if data != nil {
		t.Fatalf("%s expects data to be nil but has %v", t.Name(), data)
	} else if err != nil {
		t.Fatalf("%s expects error to be nil but has %v", t.Name(), err)
	}
}

func TestBackground_Canceled(t *testing.T) {
	flow := Sequence(
		Background(complete(nil)),
		Delay(time.Millisecond, cancel(nil)),
	)
	result, data, err := floc.Run(flow)
	if !result.IsCanceled() {
		t.Fatalf("%s expects result to be Canceled but has %s", t.Name(), result.String())
	} else if data != nil {
		t.Fatalf("%s expects data to be nil but has %v", t.Name(), data)
	} else if err != nil {
		t.Fatalf("%s expects error to be nil but has %v", t.Name(), err)
	}
}

func TestBackground_Failed(t *testing.T) {
	flow := Sequence(
		Background(cancel(nil)),
		Delay(time.Millisecond, fail(nil, fmt.Errorf("err"))),
	)
	result, data, err := floc.Run(flow)
	if !result.IsFailed() {
		t.Fatalf("%s expects result to be Failed but has %s", t.Name(), result.String())
	} else if data != nil {
		t.Fatalf("%s expects data to be nil but has %v", t.Name(), data)
	} else if err == nil {
		t.Fatalf("%s expects error to be not nil", t.Name())
	}
}

func TestBackground_Error(t *testing.T) {
	flow := Sequence(
		Background(cancel(nil)),
		Delay(time.Millisecond, throw(fmt.Errorf("err"))),
	)
	result, data, err := floc.Run(flow)
	if !result.IsFailed() {
		t.Fatalf("%s expects result to be Failed but has %s", t.Name(), result.String())
	} else if data != nil {
		t.Fatalf("%s expects data to be nil but has %v", t.Name(), data)
	} else if err == nil {
		t.Fatalf("%s expects error to be not nil", t.Name())
	}
}
