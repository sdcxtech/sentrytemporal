package sentrytemporal

import (
	"context"
	"errors"

	"github.com/getsentry/sentry-go"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/interceptor"
)

type activityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
	root *workerInterceptor
}

func (a *activityInboundInterceptor) ExecuteActivity(
	ctx context.Context,
	in *interceptor.ExecuteActivityInput,
) (ret interface{}, err error) {
	hub := a.root.hub.Clone()
	ctx = sentry.SetHubOnContext(ctx, hub)

	configureScope := func(scope *sentry.Scope) {
		info := activity.GetInfo(ctx)
		scope.SetContext("activity info", mustStruct2Map(info))
		scope.SetContext("execute activity input", mustStruct2Map(in))

		scope.SetTag("temporal_io_kind", "ExecuteActivity")

		scope.SetFingerprint(
			[]string{
				info.TaskQueue,
				info.ActivityType.Name,
				"{{ default }}",
			},
		)

		if a.root.options.ActivityScopeCustomizer != nil {
			a.root.options.ActivityScopeCustomizer(ctx, scope, err)
		}
	}

	defer func() {
		if x := recover(); x != nil {
			hub.ConfigureScope(configureScope)
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetLevel(sentry.LevelFatal)
			})
			_ = hub.Recover(x)
			panic(x)
		}
	}()

	ret, err = a.Next.ExecuteActivity(ctx, in)
	if err != nil {
		if errors.Is(err, activity.ErrResultPending) {
			return
		}

		if skipper := a.root.options.ActivityErrorSkipper; skipper != nil && skipper(ctx, err) {
			return
		}

		hub.ConfigureScope(configureScope)
		_ = hub.CaptureException(err)
	}

	return
}
