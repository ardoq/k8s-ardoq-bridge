package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalEventOps = promauto.NewCounter(prometheus.CounterOpts{
		Name: "k8sardoqbridge_event_ops_total",
		Help: "The total number of processed events",
	})
	CacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "k8sardoqbridge_cache_hit_total",
		Help: "Total Number of items fetched from the cache",
	})
	CacheMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "k8sardoqbridge_cache_miss_total",
		Help: "Total Number of items fetched from remote source",
	})
	CachePersists = promauto.NewCounter(prometheus.CounterOpts{
		Name: "k8sardoqbridge_cache_persist_total",
		Help: "Total Number of items persisted to cache",
	})
	RequestLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "k8sardoqbridge_request_latency",
		Help: "Request latencies",
	},
		[]string{"action"})
)
