package storage

import (
	"github.com/brocaar/loraserver/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	queryTimer func(string, func() error) error
)

func init() {
	qt := metrics.MustRegisterNewTimerWithError(
		"storage_function_query_duration",
		"Per internal/storage function query duration tracking.",
		[]string{"function"},
	)

	queryTimer = func(fName string, f func() error) error {
		return qt(prometheus.Labels{"function": fName}, f)
	}
}
