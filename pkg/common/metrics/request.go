package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	cm "etop.vn/backend/pkg/common"
)

var countRequests = registerCounterVec(prometheus.CounterOpts{
	Namespace:   "",
	Subsystem:   "",
	Name:        "AA_main_api_requests_total",
	Help:        "",
	ConstLabels: nil,
}, "name", "code")

var apiRequestHistogram = registerHistogramVec(prometheus.HistogramOpts{
	Namespace:   "",
	Subsystem:   "",
	Name:        "AA_main_api_requests_histogram",
	Help:        "",
	ConstLabels: nil,
	Buckets:     nil,
}, "name", "code")

func APIRequest(name string, d time.Duration, err error) {
	code := cm.ErrorCode(err).String()
	countRequests.WithLabelValues(name, code).Inc()
	apiRequestHistogram.WithLabelValues(name, code).Observe(d.Seconds())
}
