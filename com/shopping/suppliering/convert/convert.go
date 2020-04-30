package convert

import (
	"fmt"
	"time"

	"o.o/api/shopping/suppliering"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/shopping/suppliering/model  -> o.o/api/shopping/suppliering
// +gen:convert: o.o/api/shopping/suppliering

const (
	MaxCodeNorm = 999999
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
