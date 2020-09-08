package imcsv

import (
	"context"

	"o.o/api/main/catalog"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/com/main/ordering/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/imcsv"
)

// - Duplicated code
// - Code exists in database
// - Variant codes

func (im *Import) verifyOrders(ctx context.Context, shop *identitymodel.Shop, idx imcsv.Indexer, codeMode Mode, rowOrders []*RowOrder) (errs []error, _ error) {
	mapCodes := make(map[string]*RowOrder)
	codes := make([]string, len(rowOrders))
	variantCodesMap := make(map[string]*RowOrderLine)

	for i, rowOrder := range rowOrders {
		codes[i] = rowOrder.OrderEdCode
		for _, line := range rowOrder.Lines {
			if line.VariantEdCode != "" {
				variantCodesMap[line.VariantEdCode] = line
			}
		}

		prevOrder, exists := mapCodes[rowOrder.OrderEdCode]
		if !exists {
			mapCodes[rowOrder.OrderEdCode] = rowOrder
			continue
		}

		err := imcsv.CellError(idx,
			prevOrder.RowIndex, idxOrderEdCode,
			`Mã đơn hàng "%v" đã trùng lặp (mã trước ở ô %v).`,
			rowOrder.OrderEdCode,
			imcsv.CellName(prevOrder.RowIndex, idxOrderEdCode),
		)
		errs = append(errs, err)
		if len(errs) >= MaxCellErrors {
			return
		}
	}
	if len(errs) > 0 {
		return
	}
	if len(mapCodes) != len(rowOrders) {
		return nil, cm.Errorf(cm.Internal, nil, "Unexpected map length")
	}

	// Verify order codes
	orderCodeQuery := &modelx.VerifyOrdersByEdCodeQuery{
		ShopID:           shop.ID,
		EdCodes:          codes,
		OnlyActiveOrders: true,
	}
	if err := im.OrderStore.VerifyOrdersByEdCode(ctx, orderCodeQuery); err != nil {
		return nil, err
	}
	existingCodes := orderCodeQuery.Result.EdCodes

	// Shop has not created any product source yet
	if len(variantCodesMap) != 0 {
		var line *RowOrderLine
		for _, ln := range variantCodesMap {
			line = ln
			break
		}
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Cửa hàng chưa tạo sản phẩm nhưng vẫn điền mã sản phẩm. Vui lòng thêm sản phẩm vào cửa hàng hoặc xóa mã sản phẩm khỏi file import (ô %v). Nếu cần thêm thông tin vui lòng liên hệ %v.", imcsv.CellName(line.RowIndex, idxVariantEdCode), wl.X(ctx).CSEmail)
	}

	var variantMap map[string]*catalog.ShopVariantWithProduct

	// Verify variant codes
	if len(variantCodesMap) != 0 {
		variantCodes := make([]string, 0, len(variantCodesMap))
		for code := range variantCodesMap {
			variantCodes = append(variantCodes, code)
		}

		existingVariants, err := im.shopVariantStore(ctx).
			ShopID(shop.ID).
			FilterForImport(catalogsqlstore.ListShopVariantsForImportArgs{
				Codes: variantCodes,
			}).
			ListShopVariantsWithProduct()
		if err != nil {
			return nil, err
		}

		variantMap = make(map[string]*catalog.ShopVariantWithProduct)
		for _, v := range existingVariants {
			variantMap[v.Code] = v
		}
	}

	mapExistingCodes := make(map[string]struct{})
	for _, code := range existingCodes {
		mapExistingCodes[code] = struct{}{}
	}

	for _, rowOrder := range rowOrders {
		_, exists := mapExistingCodes[rowOrder.OrderEdCode]
		if exists {
			err := imcsv.CellError(idx,
				rowOrder.RowIndex, idxOrderEdCode,
				`Mã đơn hàng "%v" đã tồn tại, vui lòng sử dụng mã khác.`,
				rowOrder.OrderEdCode,
			)
			errs = append(errs, err)
			if len(errs) >= MaxCellErrors {
				return
			}
		}

		// Verify variant code/ed_code
		if variantMap == nil {
			continue
		}
		for _, line := range rowOrder.Lines {
			if line.VariantEdCode == "" {
				continue
			}

			v := variantMap[line.VariantEdCode]
			if v == nil {
				err := imcsv.CellError(idx,
					line.RowIndex, idxVariantEdCode,
					`Mã phiên bản sản phẩm "%v" không tồn tại. Vui lòng kiểm tra lại. Nếu cần thêm thông tin vui lòng liên hệ %v.`,
					line.VariantEdCode,
					wl.X(ctx).CSEmail,
				)
				errs = append(errs, err)
				if len(errs) >= MaxCellErrors {
					return
				}
				continue
			}
			line.XVariant = v
		}
	}
	return
}
