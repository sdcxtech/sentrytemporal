package workflow_cases

import (
	"errors"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

func WorkflowPanic(ctx workflow.Context, name string) (string, error) {
	panic("test workflow panic")
	// return "111", nil
}

func WorkflowError(ctx workflow.Context, name string) (string, error) {
	return "", errors.New("test workflow error")
}

func WorkflowQueryHandler(ctx workflow.Context, name string) (string, error) {
	workflow.SetQueryHandler(ctx, "error", func(queryIn string) (string, error) {
		return "", fmt.Errorf("test query handler error: %s", queryIn)
	})
	workflow.SetQueryHandler(ctx, "panic", func(queryIn string) (string, error) {
		panic(fmt.Errorf("test query handler panic: %s", queryIn))
	})

	workflow.Sleep(ctx, 30*time.Second)

	return "", nil
}

func WorkflowSignalHandler(ctx workflow.Context, name string) (string, error) {
	errorChan := workflow.GetSignalChannel(ctx, "panic")
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(errorChan, func(c workflow.ReceiveChannel, more bool) {
		panic("test signal handler panic")
	})

	selector.Select(ctx)

	return "", nil
}
