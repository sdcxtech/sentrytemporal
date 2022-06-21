package sentrytemporal

import "errors"

type continueAsNewError struct {
	wrapped error
}

func newContinueAsNewError(err error) *continueAsNewError {
	return &continueAsNewError{wrapped: err}
}

func (e *continueAsNewError) Error() string {
	return e.wrapped.Error()
}

func (e *continueAsNewError) Unwrap() error {
	return e.wrapped
}

func isContinueAsNewError(err error) bool {
	var continueAsNewErr *continueAsNewError
	return errors.As(err, &continueAsNewErr)
}
