package convert

import (
	"fmt"
	"time"

	"etop.vn/api/shopping/suppliering"
	_ "etop.vn/backend/com/shopping/suppliering/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/shopping/suppliering/model -> etop.vn/api/shopping/suppliering
// +gen:convert: etop.vn/api/shopping/suppliering

const (
	MaxCodeNorm = 999999
	codeRegex   = "^NCC([0-9]{6})$"
	codePrefix  = "NCC"
)

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func createShopSupplier(args *suppliering.CreateSupplierArgs, out *suppliering.ShopSupplier) {
	apply_suppliering_CreateSupplierArgs_suppliering_ShopSupplier(args, out)
	out.ID = cm.NewID()
}

func updateShopSupplier(args *suppliering.UpdateSupplierArgs, out *suppliering.ShopSupplier) {
	apply_suppliering_UpdateSupplierArgs_suppliering_ShopSupplier(args, out)
	out.UpdatedAt = time.Now()
}
