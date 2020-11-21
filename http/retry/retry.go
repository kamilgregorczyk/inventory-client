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

func NewRetries(config RetriesConfig) (*Retry, error) {
	if config.MaxRetries <= 0 {
		return nil, MaxRetriesBellowZeroError
	}
	if config.Delay.Milliseconds() <= 0 {
		return nil, DelayBellowZeroError
	}
	if config.Factor <= 0 {
		return nil, FactorZeroError
	}

	return &Retry{config: config}, nil
}

type Retry struct {
	config RetriesConfig
}

type RetryFunc func() (*http.Response, error)

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
