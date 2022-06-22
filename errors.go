package sentrytemporal

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

func isContinueAsNewError(err error) bool {
	var continueAsNewErr *workflow.ContinueAsNewError
	return errors.As(err, &continueAsNewErr)
}
