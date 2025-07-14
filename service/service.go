package service

import (
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sony/gobreaker"
)

type APIService interface {
	GetExternalData() (string, error)
}

type apiService struct {
	client *retryablehttp.Client
	cb     *gobreaker.CircuitBreaker
}

func NewAPIService() APIService {
	// Retryable HTTP Client
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = 500 * time.Millisecond
	client.RetryWaitMax = 2 * time.Second
	client.Backoff = retryablehttp.DefaultBackoff
	client.Logger = nil // disable log, bisa diganti custom

	// Circuit Breaker
	cbSettings := gobreaker.Settings{
		Name:    "ExternalAPI",
		Timeout: 10 * time.Second, // waktu tunggu sebelum half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}
	cb := gobreaker.NewCircuitBreaker(cbSettings)

	return &apiService{client: client, cb: cb}
}

func (s *apiService) GetExternalData() (string, error) {
	url := "https://httpstat.us/503" // simulasi API gagal

	// Bungkus dalam circuit breaker
	result, err := s.cb.Execute(func() (interface{}, error) {
		req, err := retryablehttp.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("request error: %w", err)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("retry failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read error: %w", err)
		}

		return string(body), nil
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}
