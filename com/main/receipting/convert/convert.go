package convert

import (
	"time"

	"etop.vn/api/main/receipting"
	"etop.vn/backend/com/main/receipting/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/receipting/model -> etop.vn/api/main/receipting
// +gen:convert: etop.vn/api/main/receipting

func createReceipt(args *receipting.CreateReceiptArgs, out *receipting.Receipt) {
	apply_receipting_CreateReceiptArgs_receipting_Receipt(args, out)
	out.ID = cm.NewID()
	out.Lines = args.Lines
}

func updateReceipt(args *receipting.UpdateReceiptArgs, out *receipting.Receipt) {
	apply_receipting_UpdateReceiptArgs_receipting_Receipt(args, out)
	out.Lines = args.Lines
	out.UpdatedAt = time.Now()
}

func receiptDB(args *receipting.Receipt, out *model.Receipt) {
	convert_receipting_Receipt_receiptingmodel_Receipt(args, out)
	out.OrderIDs = args.GetOrderIDs()
}
