package main

import (
	"context"
	"log"

	"github.com/sdcxtech/sentrytemporal/tests/workflow_cases"
	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	RunWorkflowPanic(c)
	RunWorkflowError(c)
	RunWorkflowQueryHandler(c)
}

func RunWorkflowPanic(c client.Client) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "sentry_tests_workflows_cases/panic",
		TaskQueue: "sentry",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow_cases.WorkflowPanic, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}

func RunWorkflowError(c client.Client) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "sentry_tests_workflows_cases/error",
		TaskQueue: "sentry",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow_cases.WorkflowError, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}

func RunWorkflowQueryHandler(c client.Client) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        "sentry_tests_workflows_cases/query",
		TaskQueue: "sentry",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, workflow_cases.WorkflowQueryHandler, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	_, err = c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), "error", "query error")
	if err != nil {
		log.Println("Unable to execute query workflow", err)
	}

	_, err = c.QueryWorkflow(context.Background(), we.GetID(), we.GetRunID(), "panic", "query panic")
	if err != nil {
		log.Println("Unable to execute query workflow", err)
	}
}
