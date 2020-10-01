package query

import (
	"context"
	"strconv"

	"o.o/api/main/shippingcode"
	cm "o.o/api/top/types/common"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/sql/cmsql"
)

var _ shippingcode.QueryService = &QueryService{}

type QueryService struct {
	db *cmsql.Database
}

func NewQueryService(db com.MainDB) *QueryService {
	return &QueryService{
		db: db,
	}
}

func QueryServiceMessageBus(q *QueryService) shippingcode.QueryBus {
	b := bus.New()
	return shippingcode.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *QueryService) GenerateShippingCode(
	ctx context.Context, args *cm.Empty,
) (string, error) {
	var code int
	if err := q.db.SQL(`SELECT nextval('shipping_code')`).Scan(&code); err != nil {
		return "", err
	}
	// checksum: avoid input wrong code

	checksumDigit := gencode.CheckSumDigitUPC(strconv.Itoa(code))
	return checksumDigit, nil
}
