package driver

import (
	"strconv"

	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/sql/cmsql"
)

type ShippingCodeGenerator interface {
	GenerateVtpostShippingCode() (string, error)
}

type generator struct {
	db *cmsql.Database
}

func NewShippingCodeGenerator(db com.MainDB) ShippingCodeGenerator {
	return generator{db: db}
}

// TODO(vu): db is only needed for generating code, use an aggregate here instead
func (g generator) GenerateVtpostShippingCode() (string, error) {
	var code int
	if err := g.db.SQL(`SELECT nextval('shipping_code')`).Scan(&code); err != nil {
		return "", err
	}
	// checksum: avoid input wrong code
	checksumDigit := gencode.CheckSumDigitUPC(strconv.Itoa(code))
	return checksumDigit, nil
}
