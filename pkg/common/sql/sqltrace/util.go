package sqltrace

import (
	"crypto/sha256"
	"fmt"

	pg_query "github.com/lfittl/pg_query_go"
)

func Fingerprint(s string) (string, bool) {
	result, err := pg_query.FastFingerprint(s)
	if err == nil {
		return result, true
	}
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("#{%x}", sum), false
}
