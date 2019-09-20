package payment

import (
	"strconv"
	"strings"

	"etop.vn/api/external/payment"
	cm "etop.vn/backend/pkg/common"
)

func ParseCode(code string) (payment.PaymentSource, int64, error) {
	// Format: order_1092423469033452748
	args := strings.Split(code, "_")
	if len(args) < 2 {
		return "", 0, cm.Errorf(cm.InvalidArgument, nil, "Wrong format")
	}
	id, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return "", 0, cm.Errorf(cm.InvalidArgument, err, "ID does not exist")
	}
	return payment.PaymentSource(args[0]), id, nil
}
