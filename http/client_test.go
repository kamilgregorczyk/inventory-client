package http

import (
	"github.com/stretchr/testify/assert"
	retry2 "test2/http/retry"
	"testing"
	"time"
)

func TestNewClientWithValidConfig(t *testing.T) {
	config := ClientConfig{
		Retries: retry2.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
		Timeout: time.Second,
		Headers: Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", config.Retries, config.Timeout, config.Headers)

	t.Logf("When creating Client")
	client, err := NewClient(config)

	t.Logf("Should not return any errors")
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientWithInValidConfig(t *testing.T) {
	testCases := []struct {
		Timeout       time.Duration
		Headers       Headers
		Retries       retry2.RetriesConfig
		ExpectedError error
	}{
		{
			Retries:       retry2.RetriesConfig{MaxRetries: 3, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Nanosecond * 0,
			Headers:       Headers{},
			ExpectedError: TimeoutZeroError,
		},
		{
			Retries:       retry2.RetriesConfig{MaxRetries: -1, Delay: time.Millisecond, Factor: 2},
			Timeout:       time.Second,
			Headers:       Headers{},
			ExpectedError: retry2.MaxRetriesZeroError,
		},
	}

	for _, testCase := range testCases {

		t.Logf("Given valid ClientConfig retries=%+v timeout=%s headers=%+v", testCase.Retries, testCase.Timeout, testCase.Headers)
		config := ClientConfig{
			Retries: testCase.Retries,
			Timeout: testCase.Timeout,
			Headers: testCase.Headers,
		}

		t.Logf("When creating Client")
		client, err := NewClient(config)

		t.Logf("Should return '%s' error", testCase.ExpectedError)
		assert.EqualError(t, err, testCase.ExpectedError.Error())
		assert.Nil(t, client)
	}
}
