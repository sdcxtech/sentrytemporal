package sentrytemporal

import (
	"context"

	"github.com/getsentry/sentry-go"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

type Options struct {
	// ActivityErrorSkipper configures a function to determine if an error from activity should be skipped.
	ActivityErrorSkipper ErrorSkipper
	// WorkflowErrorSkipper configures a function to determine if an error from workflow should be skipped.
	WorkflowErrorSkipper ErrorSkipper
}

type ErrorSkipper func(err error) bool

// New creates a worker interceptor which will report error to sentry.
func New(hub *sentry.Hub, opts Options) interceptor.WorkerInterceptor {
	i := &workerInterceptor{
		options: opts,
		hub:     hub,
	}

	if i.hub == nil {
		i.hub = sentry.CurrentHub()
	}

	return i
}

type workerInterceptor struct {
	interceptor.WorkerInterceptorBase
	options Options
	hub     *sentry.Hub
}

func (w *workerInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	i := &activityInboundInterceptor{root: w}
	i.Next = next

	return i
}

func (w *workerInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	i := &workflowInboundInterceptor{root: w}
	i.Next = next

	return i
}
