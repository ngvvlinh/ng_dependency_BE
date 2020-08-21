package inventory

import (
	"o.o/api/main/inventory"
	"o.o/api/top/int/shop"
	"o.o/backend/pkg/common/apifw/cmapi"
)

func PbInventory(args *inventory.InventoryVariant) *shop.InventoryVariant {
	return &shop.InventoryVariant{
		ShopId:         args.ShopID,
		VariantId:      args.VariantID,
		QuantityOnHand: args.QuantityOnHand,
		QuantityPicked: args.QuantityPicked,
		Quantity:       args.QuantitySummary,
		CostPrice:      args.CostPrice,
		CreatedAt:      cmapi.PbTime(args.CreatedAt),
		UpdatedAt:      cmapi.PbTime(args.UpdatedAt),
	}
}

func PbInventoryVariants(args []*inventory.InventoryVariant) []*shop.InventoryVariant {
	var inventoryVariants []*shop.InventoryVariant
	for _, value := range args {
		inventoryVariants = append(inventoryVariants, PbInventory(value))
	}
	return inventoryVariants
}

func PbShopInventoryVoucher(args *inventory.InventoryVoucher) *shop.InventoryVoucher {
	if args == nil {
		return nil
	}

	var inventoryVoucherItem []*shop.InventoryVoucherLine
	for _, value := range args.Lines {
		inventoryVoucherItem = append(inventoryVoucherItem, &shop.InventoryVoucherLine{
			VariantId:   value.VariantID,
			VariantName: value.VariantName,
			ProductId:   value.ProductID,
			Code:        value.Code,
			ProductName: value.ProductName,
			ImageUrl:    value.ImageURL,
			Attributes:  value.Attributes,
			Price:       value.Price,
			Quantity:    value.Quantity,
		})
	}
	return &shop.InventoryVoucher{
		RefAction:    args.RefAction,
		Title:        args.Title,
		TotalAmount:  args.TotalAmount,
		CreatedBy:    args.CreatedBy,
		UpdatedBy:    args.UpdatedBy,
		Lines:        inventoryVoucherItem,
		RefId:        args.RefID,
		Code:         args.Code,
		RefCode:      args.RefCode,
		RefType:      args.RefType.String(),
		RefName:      args.RefName,
		TraderId:     args.TraderID,
		Trader:       PbShopTrader(args.Trader),
		Status:       args.Status,
		Note:         args.Note,
		Type:         args.Type.String(),
		Id:           args.ID,
		ShopId:       args.ShopID,
		CreatedAt:    cmapi.PbTime(args.CreatedAt),
		UpdatedAt:    cmapi.PbTime(args.UpdatedAt),
		CancelledAt:  cmapi.PbTime(args.CancelledAt),
		ConfirmedAt:  cmapi.PbTime(args.ConfirmedAt),
		CancelReason: args.CancelReason,
	}
}

func PbShopInventoryVouchers(inventory []*inventory.InventoryVoucher) []*shop.InventoryVoucher {
	var inventoryVouchers []*shop.InventoryVoucher
	for _, value := range inventory {
		inventoryVouchers = append(inventoryVouchers, PbShopInventoryVoucher(value))
	}
	return inventoryVouchers
}

func PbShopTrader(args *inventory.Trader) *shop.Trader {
	if args == nil {
		return nil
	}
	return &shop.Trader{
		Id:       args.ID,
		Type:     args.Type,
		FullName: args.FullName,
		Phone:    args.Phone,
		Deleted:  false,
	}
}
