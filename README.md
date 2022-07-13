# Sentry Temporal Interceptor

The [temporal](https://www.temporal.io/) interceptor captures panic and errors and report them 
to [sentry](https://github.com/getsentry/sentry) server.

## Installation

```sh
go get github.com/sdcxtech/sentrytemporal
```

```go
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "taskQueue", worker.Options{
		// Create interceptor that will put started time on the logger
		Interceptors: []sdkinterceptor.WorkerInterceptor{
			sentrytemp.New(
				sentry.CurrentHub(),
				sentrytemp.Options{},
			),
		},
	})
```

## Configuration

`New` accepts a struct of Options that allows you to configure how the inteceptor will behave.

```go
type Options struct {
	// ActivityErrorSkipper configures a function to determine if an error from activity should be skipped.
	ActivityErrorSkipper ActivityErrorSkipper
	// WorkflowErrorSkipper configures a function to determine if an error from workflow should be skipped.
	WorkflowErrorSkipper WorkflowErrorSkipper
}

type (
	ActivityErrorSkipper func(context.Context, error) bool
	WorkflowErrorSkipper func(workflow.Context, error) bool
)
```

Example:

Only report retrieable error when attempt count is great then `1`.

```go
activityErrorSkipper := func(ctx context.Context, err error) bool {
	var errApp *temporal.ApplicationError
	if errors.As(err, &errApp) {
		if errApp.NonRetryable() {
			return false
		}

		info := activity.GetInfo(ctx)

		return !(info.Attempt > 1)
	}

	return false
}
```
