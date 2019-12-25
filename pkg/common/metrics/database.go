package metrics

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"etop.vn/backend/pkg/common/sql/sq"
)

var databaseQuery = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "AA_main_database_query_histogram",
}, "tx", "query")

func DatabaseQuery(fingerprint string, entry *sq.LogEntry) {
	databaseQuery.WithLabelValues(
		strconv.FormatBool(entry.IsTx()),
		fingerprint,
	).Observe(entry.Duration.Seconds())
}

var databaseTransaction = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "AA_main_database_transaction_histogram",
}, "n", "type", "query")

func DatabaseTransaction(fingerprints []string, entry *sq.LogEntry) {
	flag, _ := entry.Flags.MarshalJSON()
	databaseTransaction.WithLabelValues(
		strconv.Itoa(len(entry.TxQueries)),
		string(flag),
		strings.Join(fingerprints, ","),
	).Observe(entry.Duration.Seconds())
}
