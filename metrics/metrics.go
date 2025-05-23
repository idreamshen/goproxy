package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CounterProxyProcessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "app_proxy_processed_total",
	})

	HistogramProxyProcessSec = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_proxy_processed_seconds",
		},
		[]string{"host"},
	)
)
