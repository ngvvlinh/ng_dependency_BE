package convert

import (
	"fmt"
	"time"

	"o.o/api/main/inventory"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/backend/com/main/inventory/model"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/inventory/model -> o.o/api/main/inventory
// +gen:convert: o.o/api/main/inventory

const (
	MaxCodeNorm = 999999
	codePrefix  = "PTK"
)

func GenerateCode(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func createInventoryVoucher(args *inventory.CreateInventoryVoucherArgs, out *inventory.InventoryVoucher) {
	apply_inventory_CreateInventoryVoucherArgs_inventory_InventoryVoucher(args, out)
	out.ID = cm.NewID()
	out.UpdatedBy = out.CreatedBy
	if out.Title == "" {
		if args.Type == inventory_type.In {
			out.Title = "Phiếu nhập kho"
		} else if args.Type == inventory_type.Out {
			out.Title = "Phiếu xuất kho"
		}
	}
}

func inventoryVariantModel(args *inventory.InventoryVariant,
	out *model.InventoryVariant) *model.InventoryVariant {
	if out == nil {
		out = &model.InventoryVariant{}
	}
	convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(args, out)
	return out
}

func inventoryVariant(args *model.InventoryVariant,
	out *inventory.InventoryVariant) *inventory.InventoryVariant {
	convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(args, out)
	out.QuantitySummary = args.QuantityOnHand + args.QuantityPicked
	return out
}

func updateInventoryVoucher(arg *inventory.UpdateInventoryVoucherArgs,
	out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	apply_inventory_UpdateInventoryVoucherArgs_inventory_InventoryVoucher(arg, out)
	out.UpdatedAt = time.Time{}
	return out
}

func inventoryVoucherModel(args *inventory.InventoryVoucher,
	out *model.InventoryVoucher) *model.InventoryVoucher {
	if args == nil {
		return nil
	}
	if out == nil {
		out = &model.InventoryVoucher{}
	}
	convert_inventory_InventoryVoucher_inventorymodel_InventoryVoucher(args, out)
	return out
}

func inventoryVoucher(args *model.InventoryVoucher,
	out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	if args == nil {
		return nil
	}
	if out == nil {
		out = &inventory.InventoryVoucher{}
	}
	convert_inventorymodel_InventoryVoucher_inventory_InventoryVoucher(args, out)
	return out
}
