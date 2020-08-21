package imcsv

import (
	"context"
	"strings"

	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/validate"
)

func (im *Import) verifyFulfillments(ctx context.Context, shop *identitymodel.Shop, idx imcsv.Indexer, rowFulfillments []*RowFulfillment) (errs []error, _ error) {
	mapCodes := make(map[string]*RowFulfillment)
	codes := make([]string, len(rowFulfillments))

	for i, rowFulfillment := range rowFulfillments {
		// check total_weight >= 50 gram
		if rowFulfillment.TotalWeight < 50 {
			err := imcsv.CellError(idx,
				rowFulfillment.RowIndex, idxTotalWeight,
				"Khối lương đơn giao hàng phải lớn hơn hoặc bằng 50 gram")
			errs = append(errs, err)
		}

		// check cod_amount >= 5000 dong
		if rowFulfillment.CODAmount != 0 && rowFulfillment.CODAmount < 5000 {
			err := imcsv.CellError(idx,
				rowFulfillment.RowIndex, idxCODAmount,
				"Tiền thu hộ phải bằng 0 hoặc từ 5000 đồng trở lên")
			errs = append(errs, err)
		}

		// check phone
		phoneNorm, ok := validate.NormalizePhone(rowFulfillment.CustomerPhone)
		if !ok {
			err := imcsv.CellError(idx,
				rowFulfillment.RowIndex, idxCustomerPhone,
				"Số điện thoại khách hàng không hợp lệ")
			errs = append(errs, err)
		} else {
			rowFulfillment.CustomerPhone = phoneNorm.String()
		}

		if strings.TrimSpace(rowFulfillment.EdCode) == "" {
			continue
		}
		codes[i] = rowFulfillment.EdCode

		prevFulfillment, exists := mapCodes[rowFulfillment.EdCode]
		if !exists {
			mapCodes[rowFulfillment.EdCode] = rowFulfillment
			continue
		}

		err := imcsv.CellError(idx,
			prevFulfillment.RowIndex, idxFulfillmentEdCode,
			`Mã đơn hàng "%v" đã trùng lặp (mã trước ở ô %v).`,
			rowFulfillment.EdCode,
			imcsv.CellName(prevFulfillment.RowIndex, idxFulfillmentEdCode),
		)
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return
	}

	// TODO(ngoc): Verify fulfillmentEdCode duplicate in database

	return
}
