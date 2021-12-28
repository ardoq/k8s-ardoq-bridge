package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalEventOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ardoqk8sbridge_event_ops_total",
		Help: "The total number of processed events",
	})
)
