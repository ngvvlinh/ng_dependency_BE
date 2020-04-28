package metrics

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"o.o/backend/pkg/common/sql/sq"
)

var databaseQuery = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "AA_main_database_query_histogram",
}, "tx", "query", "error")

func DatabaseQuery(fingerprint string, entry *sq.LogEntry) {
	databaseQuery.WithLabelValues(
		strconv.FormatBool(entry.IsTx()),
		fingerprint,
		errorStr(entry.OrigError),
	).Observe(entry.Duration.Seconds())
}

var databaseTransaction = registerHistogramVec(prometheus.HistogramOpts{
	Namespace: "",
	Subsystem: "",
	Name:      "AA_main_database_transaction_histogram",
}, "n", "type", "query", "error")

func DatabaseTransaction(fingerprints []string, entry *sq.LogEntry) {
	flag, _ := entry.Flags.MarshalJSON()
	databaseTransaction.WithLabelValues(
		strconv.Itoa(len(entry.TxQueries)),
		string(flag),
		strings.Join(fingerprints, ","),
		errorStr(entry.OrigError),
	).Observe(entry.Duration.Seconds())
}

func errorStr(err error) string {
	if err == nil {
		return "ok"
	}
	return err.Error()
}
