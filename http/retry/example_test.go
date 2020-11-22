package retry_test

import (
	"log"
	"net/http"
	retries "test2/http/retry"
	"time"
)

func ExampleRetry_Execute() {
	config := retries.RetriesConfig{
		MaxRetries: 3,
		Delay:      time.Millisecond * 500,
		Factor:     1.3,
	}

	retry, err := retries.NewRetries(config)

	if err != nil {
		log.Fatal(err)
	}
	retry.Execute(func() (*http.Response, error) {
		response, err := http.Get("http://localhost")
		// We need retries only on 500s and higher
		if response.StatusCode >= 500 {
			return response, &retries.RetryableError{Err: err}
		}
		return response, err
	})
}

func ExampleNewRetries() {
	client, err := retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 3,
		Delay:      time.Millisecond * 500,
		Factor:     1.3,
	})

	client, err = retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 5,
		Delay:      time.Second,
		Factor:     2,
	})

	client, err = retries.NewRetries(retries.RetriesConfig{
		MaxRetries: 1,
		Delay:      time.Second,
		Factor:     1,
	})

	log.Print(client, err)
}
