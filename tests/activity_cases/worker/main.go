package main

import (
	"log"

	"github.com/getsentry/sentry-go"
	"go.temporal.io/sdk/client"
	sdkinterceptor "go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"

	sentrytemp "github.com/sdcxtech/sentrytemporal"
	"github.com/sdcxtech/sentrytemporal/tests/activity_cases"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatalln("init sentry failed", err)
	}

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "sentry", worker.Options{
		// Create interceptor that will put started time on the logger
		Interceptors: []sdkinterceptor.WorkerInterceptor{
			sentrytemp.New(
				sentry.CurrentHub(),
				sentrytemp.Options{},
			),
		},
	})

	w.RegisterWorkflow(activity_cases.Workflow)
	w.RegisterActivity(activity_cases.ActivityPanic)
	w.RegisterActivity(activity_cases.ActivityError)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
