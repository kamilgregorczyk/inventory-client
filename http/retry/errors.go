package retry

import "errors"

var (
	MaxRetriesZeroError = errors.New("maxRetries has to be larger than 0")
	DelayZeroError      = errors.New("delay has to be larger than 0")
	FactorZeroError     = errors.New("factor has to be larger than 0")
)


type RetryableError struct {
	Err error
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func (e *RetryableError) Error() string {
	if e.Err == nil {
		return "retryable: <nil>"
	}
	return "retryable: " + e.Err.Error()
}
