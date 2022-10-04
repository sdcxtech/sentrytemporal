package sentrytemporal

import (
	"context"

	"github.com/getsentry/sentry-go"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

type Options struct {
	// ActivityErrorSkipper configures a function to determine if an error from activity should be skipped.
	// If it returns true, the error is ignored.
	ActivityErrorSkipper ActivityErrorSkipper
	// WorkflowErrorSkipper configures a function to determine if an error from workflow should be skipped.
	// If it returns true, the error is ignored.
	WorkflowErrorSkipper WorkflowErrorSkipper

	// ActivityScopeCustomizer applies custom options to a sentry.Scope just before an error is reported from an activity
	ActivityScopeCustomizer ActivityScopeCustomizer
	// WorkflowScopeCustomizer applies custom options to a sentry.Scope just before an error is reported from a workflow
	WorkflowScopeCustomizer WorkflowScopeCustomizer
}

type (
	ActivityErrorSkipper func(context.Context, error) bool
	WorkflowErrorSkipper func(workflow.Context, error) bool

	ActivityScopeCustomizer func(context.Context, *sentry.Scope, error)
	WorkflowScopeCustomizer func(workflow.Context, *sentry.Scope, error)
)

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
