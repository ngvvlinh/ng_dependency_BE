package util

import (
	"fmt"

	"o.o/api/main/inventory"
	"o.o/api/main/stocktaking"
	"o.o/api/top/types/etc/inventory_type"
	"o.o/api/top/types/etc/inventory_voucher_ref"
	"o.o/api/top/types/etc/ref_action"
	"o.o/api/top/types/etc/stocktake_type"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

func PopulateInventoryVouchers(args []*inventory.InventoryVoucher, mapStocktake map[dot.ID]stocktake_type.StocktakeType) ([]*inventory.InventoryVoucher, error) {
	for key, inventoryVoucher := range args {
		if args[key].RefType == inventory_voucher_ref.StockTake {
			if mapStocktake[inventoryVoucher.RefID] == stocktake_type.Balance {
				args[key].RefName = "Kiểm kho"
				if args[key].Type == inventory_type.In {
					args[key].Note = fmt.Sprintf("Tạo phiếu nhập kho theo phiếu kiểm kho %v", args[key].RefCode)
				} else {
					args[key].Note = fmt.Sprintf("Tạo phiếu xuất kho theo phiếu kiểm kho %v", args[key].RefCode)
				}
			} else {
				args[key].RefName = "Xuất hủy"
				args[key].Note = fmt.Sprintf("Tạo phiếu xuất kho theo phiếu xuất hủy %v", args[key].RefCode)
			}
			args[key].RefAction = ref_action.Create
		} else {
			args[key] = addInfoInventoryVoucher(args[key])
		}
	}
	return args, nil
}

func PopulateInventoryVoucher(arg *inventory.InventoryVoucher, stocktake *stocktaking.ShopStocktake) (*inventory.InventoryVoucher, error) {
	if arg.RefType == inventory_voucher_ref.StockTake {
		if stocktake == nil {
			return nil, cm.Errorf(cm.Internal, nil, "Phiếu kiểm kho không tồn tại")
		}
		if stocktake.Type == stocktake_type.Balance {
			arg.RefName = "Kiểm kho"
			if arg.Type == inventory_type.In {
				arg.Note = fmt.Sprintf("Tạo phiếu nhập kho theo phiếu kiểm kho %v", arg.RefCode)
			} else {
				arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo phiếu kiểm kho %v", arg.RefCode)
			}
		} else {
			arg.RefName = "Xuất hủy"
			arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo phiếu xuất hủy %v", arg.RefCode)
		}
		arg.RefAction = ref_action.Create
	} else {
		arg = addInfoInventoryVoucher(arg)
	}
	return arg, nil
}

func addInfoInventoryVoucher(arg *inventory.InventoryVoucher) *inventory.InventoryVoucher {
	switch arg.RefType {
	case inventory_voucher_ref.Order:
		if arg.Type == inventory_type.Out {
			arg.RefAction = ref_action.Create
			arg.RefName = "Bán hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo đơn hàng %v", arg.Code)
		} else {
			arg.RefAction = ref_action.Cancel
			arg.RefName = "Hủy bán hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn hàng %v", arg.Code)
		}
	case inventory_voucher_ref.Refund:
		if arg.Type == inventory_type.In {
			arg.RefAction = ref_action.Create
			arg.RefName = "Trả hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn trả hàng  %v", arg.Code)
		} else {
			arg.RefAction = ref_action.Cancel
			arg.RefName = "Hủy trả hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo đơn trả hàng %v", arg.Code)
		}
	case inventory_voucher_ref.PurchaseOrder:
		if arg.Type == inventory_type.In {
			arg.RefAction = ref_action.Create
			arg.RefName = "Nhập hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn nhập hàng %v", arg.Code)
		} else {
			arg.RefAction = ref_action.Cancel
			arg.RefName = "Hủy nhập hàng"
			arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo đơn nhập hàng %v", arg.Code)
		}
	case inventory_voucher_ref.PurchaseRefund:
		if arg.Type == inventory_type.Out {
			arg.RefAction = ref_action.Create
			arg.RefName = "Trả hàng nhập"
			arg.Note = fmt.Sprintf("Tạo phiếu xuất kho theo đơn trả hàng nhập %v", arg.Code)
		} else {
			arg.RefAction = ref_action.Cancel
			arg.RefName = "Hủy trả hàng nhập"
			arg.Note = fmt.Sprintf("Tạo phiếu nhập kho theo đơn trả hàng nhập %v", arg.Code)
		}
	}
	return arg
}
