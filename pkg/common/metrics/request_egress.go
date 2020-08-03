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

var faboEgressRequestHistogram = registerHistogramVec(prometheus.HistogramOpts{
	Namespace:   "",
	Subsystem:   "",
	Name:        "AA_main_fabo_egress_request_histogram",
	Help:        "",
	ConstLabels: nil,
	Buckets:     nil,
}, "host", "path", "status_code", "source", "page_id")

func FaboEgressRequest(httpUrl *url.URL, httpStatusCode int, d time.Duration, source, pageID string) {
	faboEgressRequestHistogram.WithLabelValues(httpUrl.Host, httpUrl.Path, strconv.Itoa(httpStatusCode), source, pageID).Observe(d.Seconds())
}
