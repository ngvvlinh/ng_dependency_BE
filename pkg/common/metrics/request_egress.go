package metrics

import (
	"net/url"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var egressRequestHistogram = registerHistogramVec(prometheus.HistogramOpts{
	Namespace:   "",
	Subsystem:   "",
	Name:        "AA_egress_requests_histogram",
	Help:        "",
	ConstLabels: nil,
	Buckets:     nil,
}, "host", "path", "status_code")

func EgressRequest(httpUrl *url.URL, httpStatusCode int, d time.Duration) {
	egressRequestHistogram.WithLabelValues(httpUrl.Host, httpUrl.Path, strconv.Itoa(httpStatusCode)).Observe(d.Seconds())
}
