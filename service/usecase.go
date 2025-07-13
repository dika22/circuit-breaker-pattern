package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type APIService interface {
	GetExternalData() (string, error)
}

type apiService struct {
	cb *gobreaker.CircuitBreaker
}

func NewAPIService() APIService {
	settings := gobreaker.Settings{
		Name:    "ThirdPartyAPI",
		Timeout: 5 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	return &apiService{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

func (s *apiService) GetExternalData() (string, error) {
	result, err := s.cb.Execute(func() (interface{}, error) {
		// Contoh third-party API yang mungkin error
		resp, err := http.Get("https://httpstat.us/503")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return nil, fmt.Errorf("server error: %d", resp.StatusCode)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		return string(body), nil
	})

	if err != nil {
		return "", err
	}

	return result.(string), nil
}
