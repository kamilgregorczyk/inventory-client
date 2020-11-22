package retry

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewRetriesWithValidConfig(t *testing.T) {
	config := RetriesConfig{
		MaxRetries: 1,
		Delay:      time.Second,
		Factor:     1.0,
	}
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", config.MaxRetries, config.Delay, config.Factor)

	t.Logf("When creating Retry")
	retry, err := NewRetries(config)

	t.Logf("Should not return any errors")

	assert.Nil(t, err)
	assert.NotNil(t, retry)
}

func TestNewRetriesWithInValidConfig(t *testing.T) {
	configs := []struct {
		MaxRetries    int
		Delay         time.Duration
		Factor        float64
		ExpectedError error
	}{
		{MaxRetries: 0, Delay: time.Second, Factor: 1.0, ExpectedError: MaxRetriesZeroError},
		{MaxRetries: -1, Delay: time.Second, Factor: 1.0, ExpectedError: MaxRetriesZeroError},
		{MaxRetries: 1, Delay: 0 * time.Second, Factor: 1.0, ExpectedError: DelayZeroError},
		{MaxRetries: 1, Delay: -1 * time.Second, Factor: 1.0, ExpectedError: DelayZeroError},
		{MaxRetries: 1, Delay: time.Second, Factor: 0, ExpectedError: FactorZeroError},
		{MaxRetries: 1, Delay: time.Second, Factor: -1.0, ExpectedError: FactorZeroError},
	}
	for _, test := range configs {
		t.Logf("Given invalid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", test.MaxRetries, test.Delay, test.Factor)
		config := RetriesConfig{
			MaxRetries: test.MaxRetries,
			Delay:      test.Delay,
			Factor:     test.Factor,
		}

		t.Logf("When creating Retry")
		retry, err := NewRetries(config)

		t.Logf("Should return '%s' error", test.ExpectedError)
		assert.EqualError(t, err, test.ExpectedError.Error())
		assert.Nil(t, retry)
	}
}

func TestRetryWithSuccessAtFirstTry(t *testing.T) {
	maxRetries := 3
	delay := time.Millisecond
	factor := 1.0
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
	config := RetriesConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
		Factor:     factor,
	}
	t.Logf("And given Retry")
	retry, _ := NewRetries(config)

	t.Logf("And given a func to run")
	var callCount int
	expectedResponse := http.Response{}
	funcToRetry := func() (*http.Response, error) {
		callCount++
		return &expectedResponse, nil
	}

	t.Logf("When executing a func")
	response, err := retry.Execute(funcToRetry)

	t.Logf("Should call only once and not return any errors")
	assert.Equal(t, callCount, 1)
	assert.Nil(t, err)
	assert.Equal(t, &expectedResponse, response)
}

func TestRetryWithInitialFailuresAndThenSuccess(t *testing.T) {
	numberOfRetries := []int{1, 2, 3}
	for _, retryCount := range numberOfRetries {
		maxRetries := 3
		delay := time.Millisecond
		factor := 1.0
		t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
		config := RetriesConfig{
			MaxRetries: maxRetries,
			Delay:      delay,
			Factor:     factor,
		}
		t.Logf("And given Retry")
		retry, _ := NewRetries(config)

		t.Logf("And given a func to run")
		var callCount int
		expectedResponse := http.Response{}
		expectedCallCount := retryCount + 1
		funcToRetry := func() (*http.Response, error) {
			if retryCount > callCount {
				callCount++
				return &expectedResponse, &RetryableError{}
			}
			callCount++
			return &expectedResponse, nil
		}

		t.Logf("When executing a func")
		response, err := retry.Execute(funcToRetry)

		t.Logf("Should call function %d times and not return any errors", expectedCallCount)
		assert.Equal(t, callCount, expectedCallCount)
		assert.Nil(t, err)
		assert.Equal(t, &expectedResponse, response)
	}

}

func TestRetryWithConstantFailures(t *testing.T) {
	maxRetries := 3
	delay := time.Millisecond
	factor := 1.0
	t.Logf("Given valid RetriesConfig maxRetries=%d delay=%s factor=%0.2f", maxRetries, delay, factor)
	config := RetriesConfig{
		MaxRetries: maxRetries,
		Delay:      delay,
		Factor:     factor,
	}
	t.Logf("And given Retry")
	retry, _ := NewRetries(config)

	t.Logf("And given a func to run")
	var callCount int
	expectedResponse := http.Response{}
	funcToRetry := func() (*http.Response, error) {
		callCount++
		return &expectedResponse, &RetryableError{}
	}

	t.Logf("When executing a func")
	response, err := retry.Execute(funcToRetry)

	t.Logf("Should call function %d times and return errors", 4)
	if callCount != 4 {
		t.Errorf("Func should be called %d was called %d times", 4, callCount)
	}
	assert.Equal(t, callCount, 4)
	assert.NotNil(t, err)
	assert.Equal(t, &expectedResponse, response)

}
