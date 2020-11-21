package retry

import (
	"errors"
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

	if err != nil {
		t.Errorf("Failed to create Retry")
	}

	if retry == nil {
		t.Error("NewRetries didn't return Retry")
	}
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
		if !errors.Is(err, test.ExpectedError) {
			t.Errorf("Error returned is %s", err)
		}
		if retry != nil {
			t.Errorf("Retry should not be returned when error occurs")

		}
	}
}
