package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"etop.vn/backend/com/main/purchaseorder/model"

	"etop.vn/api/main/purchaseorder"
	_ "etop.vn/backend/com/main/purchaseorder/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/l"
)

// +gen:convert: etop.vn/backend/com/main/purchaseorder/model -> etop.vn/api/main/purchaseorder
// +gen:convert: etop.vn/api/main/purchaseorder

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
	out.VariantIDs = args.GetVariantIDs()
}
