package retry

import "errors"

var (
	MaxRetriesZeroError = errors.New("maxRetries has to be larger than 0")
	DelayZeroError      = errors.New("delay has to be larger than 0")
	FactorZeroError     = errors.New("factor has to be larger than 0")
)

func RetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &retryableError{err}
}

type retryableError struct {
	Err error
}

func (e *retryableError) Unwrap() error {
	return e.Err
}

func (e *retryableError) Error() string {
	if e.Err == nil {
		return "retryable: <nil>"
	}
	return "retryable: " + e.Err.Error()
}
