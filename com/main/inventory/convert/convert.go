package convert

import (
	"fmt"
	"time"

	"etop.vn/api/main/catalog"
	"etop.vn/api/main/inventory"
	"etop.vn/api/main/stocktaking"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	"etop.vn/backend/com/main/inventory/model"
	cm "etop.vn/backend/pkg/common"
)

// +gen:convert: etop.vn/backend/com/main/inventory/model->etop.vn/api/main/inventory
// +gen:convert: etop.vn/api/main/inventory

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
		if args.Type == "in" {
			out.Title = "Phiếu nhập kho"
		} else if args.Type == "out" {
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

func ConvertAttributesOrder(attributes []*catalogmodel.ProductAttribute) []*inventory.Attribute {
	var result []*inventory.Attribute
	for _, value := range attributes {
		result = append(result, &inventory.Attribute{
			Name:  value.Name,
			Value: value.Value,
		})
	}
	return result
}

func ConvertAttributesPurchaseOrder(attributes []*catalog.Attribute) []*inventory.Attribute {
	var result []*inventory.Attribute
	for _, value := range attributes {
		result = append(result, &inventory.Attribute{
			Name:  value.Name,
			Value: value.Value,
		})
	}
	return result
}

func ConvertAttributesStocktake(attributes []*stocktaking.Attribute) []*inventory.Attribute {
	var result []*inventory.Attribute
	for _, value := range attributes {
		result = append(result, &inventory.Attribute{
			Name:  value.Name,
			Value: value.Value,
		})
	}
	return result
}
