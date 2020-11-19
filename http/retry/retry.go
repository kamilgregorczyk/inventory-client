package retry

import (
	"errors"
	"math"
	"net/http"
	"time"
)

type RetriesConfig struct {
	MaxRetries int
	Delay      time.Duration
	Factor     float64
}

func NewRetries(config RetriesConfig) *Retry {
	return &Retry{config: config}
}

type Retry struct {
	config RetriesConfig
}

type RetryFunc func() (*http.Response, error)

func RetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &retryableError{err}
}

type retryableError struct {
	err error
}

func (e *retryableError) Unwrap() error {
	return e.err
}

func (e *retryableError) Error() string {
	if e.err == nil {
		return "retryable: <nil>"
	}
	return "retryable: " + e.err.Error()
}

func (r *Retry) Execute(runnable RetryFunc) (*http.Response, error) {
	var tryCount int
	for {
		response, err := runnable()
		if err == nil || tryCount >= r.config.MaxRetries {
			return response, err
		}

		var retryError *retryableError
		if !errors.As(err, &retryError) {
			return response, err
		}

		tryCount++
		select {
		case <-time.After(r.next(tryCount)):
		}

	}

}

func (r *Retry) next(currentTry int) time.Duration {
	delay := r.config.Delay.Nanoseconds() * int64(math.Pow(r.config.Factor, float64(currentTry)))
	return time.Duration(delay)

}
