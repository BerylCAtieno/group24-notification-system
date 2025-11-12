package circuit_breaker

import (
	"time"

	"github.com/sony/gobreaker/v2"
)

func NewEmailCircuitBreaker() *gobreaker.CircuitBreaker[any] {
	settings := gobreaker.Settings{
		Name:        "email_smtp_breaker",
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip when failures exceed 60% of total
			failRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests > 10 && failRatio > 0.6
		},
	}
	return gobreaker.NewCircuitBreaker[any](settings)
}
