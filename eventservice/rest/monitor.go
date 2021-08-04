package rest

import (
	"github.com/prometheus/client_golang/prometheus"
)

var eventsAddDelay = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:      "events_add_delay",
		Namespace: "myevents",
		Help:      "Delay of events created",
		Buckets:   []float64{500, 800, 1000, 2000},
	},
)

func init() {
	prometheus.Register(eventsAddDelay)
}
