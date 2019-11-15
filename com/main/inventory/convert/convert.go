package convert

import (
	"fmt"
	"time"

	"etop.vn/api/main/inventory"
	"etop.vn/backend/com/main/inventory/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/inventory/model->etop.vn/api/main/inventory
// +gen:convert: etop.vn/api/main/inventory

const (
	MaxCodeNorm = 999999
	codeRegex   = "^PTK([0-9]{6})$"
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
		if args.Type == "in" {
			out.Title = "Phiếu nhập kho"
		} else if args.Type == "out" {
			out.Title = "Phiếu xuất kho"
		}
	}
}

func InventoryVariantsFromModel(args []*model.InventoryVariant) (outs []*inventory.InventoryVariant) {
	return Convert_inventorymodel_InventoryVariants_inventory_InventoryVariants(args)
}

func InventoryVariantToModel(args *inventory.InventoryVariant,
	out *model.InventoryVariant) *model.InventoryVariant {
	if out == nil {
		out = &model.InventoryVariant{}
	}
	convert_inventory_InventoryVariant_inventorymodel_InventoryVariant(args, out)
	return out
}

func InventoryVariantFromModel(args *model.InventoryVariant,
	out *inventory.InventoryVariant) *inventory.InventoryVariant {
	if args == nil {
		return nil
	}
	if out == nil {
		out = &inventory.InventoryVariant{}
	}
	convert_inventorymodel_InventoryVariant_inventory_InventoryVariant(args, out)
	out.QuantitySummary = args.QuantityOnHand + args.QuantityPicked
	return out
}

func InventoryVouchersFromModel(args []*model.InventoryVoucher) []*inventory.InventoryVoucher {
	return Convert_inventorymodel_InventoryVouchers_inventory_InventoryVouchers(args)
}

func ApplyUpdateInventoryVoucher(arg *inventory.UpdateInventoryVoucherArgs,
	out *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	apply_inventory_UpdateInventoryVoucherArgs_inventory_InventoryVoucher(arg, out)
	out.UpdatedAt = time.Time{}
	return out
}

func InventoryVoucherToModel(args *inventory.InventoryVoucher,
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

func InventoryVoucherFromModel(args *model.InventoryVoucher,
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
