package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"o.o/api/main/purchaseorder"
	"o.o/backend/com/main/purchaseorder/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

// +gen:convert: o.o/backend/com/main/purchaseorder/model  -> o.o/api/main/purchaseorder
// +gen:convert: o.o/api/main/purchaseorder

var ll = l.New()

const (
	MaxCodeNorm = 999999
	codeRegex   = "^DNH([0-9]{6})$"
	codePrefix  = "DNH"
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

func createPurchaseOrder(args *purchaseorder.CreatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) {
	apply_purchaseorder_CreatePurchaseOrderArgs_purchaseorder_PurchaseOrder(args, out)
	out.ID = cm.NewID()
}

func updatePurchaseOrder(args *purchaseorder.UpdatePurchaseOrderArgs, out *purchaseorder.PurchaseOrder) {
	if len(args.Lines) == 0 {
		args.Lines = out.Lines
	}
	apply_purchaseorder_UpdatePurchaseOrderArgs_purchaseorder_PurchaseOrder(args, out)
	out.UpdatedAt = time.Now()
}

func purchaseOrderDB(args *purchaseorder.PurchaseOrder, out *model.PurchaseOrder) {
	convert_purchaseorder_PurchaseOrder_purchaseordermodel_PurchaseOrder(args, out)
	if out.Supplier != nil {
		out.SupplierFullNameNorm = validate.NormalizeSearch(out.Supplier.FullName)
		out.SupplierPhoneNorm = validate.NormalizeSearchPhone(out.Supplier.Phone)
	}
	out.VariantIDs = args.GetVariantIDs()
}

func purchaseOrder(args *model.PurchaseOrder, out *purchaseorder.PurchaseOrder) {
	convert_purchaseordermodel_PurchaseOrder_purchaseorder_PurchaseOrder(args, out)
	out.CancelReason = args.CancelledReason
}
