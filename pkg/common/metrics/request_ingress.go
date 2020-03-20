package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	cm "etop.vn/backend/pkg/common"
)

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
	apiRequestHistogram.WithLabelValues(name, code).Observe(d.Seconds())
}
