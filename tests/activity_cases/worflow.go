package activity_cases

import (
	"context"
	"errors"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Workflow(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string

	_ = workflow.ExecuteActivity(ctx, ActivityPanic, name).Get(ctx, &result)

	_ = workflow.ExecuteActivity(ctx, ActivityError, name).Get(ctx, &result)

	return result, nil
}

func ActivityPanic(ctx context.Context, name string) (string, error) {
	panic("test panic activty")
}

func ActivityError(ctx context.Context, name string) (string, error) {
	return "", errors.New("test actvity error")
}
