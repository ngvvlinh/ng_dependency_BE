package payment

import (
	"strconv"
	"strings"

	"o.o/api/top/types/etc/payment_source"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func ParsePaymentCode(code string) (payment_source.PaymentSource, dot.ID, error) {
	// Format: order_1092423469033452748
	parts := strings.Split(code, "_")
	if len(parts) < 2 {
		return 0, 0, cm.Errorf(cm.InvalidArgument, nil, "Wrong format")
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, cm.Errorf(cm.InvalidArgument, err, "ID does not exist")
	}
	source, ok := payment_source.ParsePaymentSource(parts[0])
	if !ok {
		return 0, 0, cm.Errorf(cm.InvalidArgument, nil, "Invalid payment source (%v)", parts[0])
	}
	return source, dot.ID(id), nil
}
