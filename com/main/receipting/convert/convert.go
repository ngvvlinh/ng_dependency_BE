package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"etop.vn/api/main/receipting"
	"etop.vn/backend/com/main/receipting/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/l"
)

// +gen:convert: etop.vn/backend/com/main/receipting/model -> etop.vn/api/main/receipting
// +gen:convert: etop.vn/api/main/receipting

var ll = l.New()

const (
	MaxCodeNorm = 999999
	codeRegex   = "^PTC([0-9]{6})$"
	codePrefix  = "PTC"
)

var reCode = regexp.MustCompile(codeRegex)

func ParseCodeNorm(code string) (_ int, ok bool) {
	parts := reCode.FindStringSubmatch(code)
	if len(parts) == 0 {
		return 0, false
	}
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		ll.Panic("unexpected", l.Error(err))
	}
	return number, true
}

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func createReceipt(args *receipting.CreateReceiptArgs, out *receipting.Receipt) {
	apply_receipting_CreateReceiptArgs_receipting_Receipt(args, out)
	out.ID = cm.NewID()
}

func updateReceipt(args *receipting.UpdateReceiptArgs, out *receipting.Receipt) {
	if len(args.Lines) == 0 {
		args.Lines = out.Lines
	}
	if args.PaidAt.IsZero() {
		args.PaidAt = out.PaidAt
	}
	apply_receipting_UpdateReceiptArgs_receipting_Receipt(args, out)
	out.UpdatedAt = time.Now()
}

func receiptDB(args *receipting.Receipt, out *model.Receipt) {
	convert_receipting_Receipt_receiptingmodel_Receipt(args, out)
	out.RefIDs = args.GetRefIDs()
	out.CreatedType = args.Mode
}

func receipt(args *model.Receipt, out *receipting.Receipt) {
	convert_receiptingmodel_Receipt_receipting_Receipt(args, out)
	out.Mode = args.CreatedType
	out.CancelReason = args.CancelledReason
}
