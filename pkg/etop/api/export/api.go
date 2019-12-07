package export

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"etop.vn/api/top/types/etc/status4"

	apishop "etop.vn/api/top/int/shop"
	ordering "etop.vn/backend/com/main/ordering/modelx"
	shipping "etop.vn/backend/com/main/shipping/modelx"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmapi"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
)

var ServiceImpl = &Service{}

type Service struct{}

func (s *Service) RequestExport(ctx context.Context, claim claims.ShopClaim, shop *model.Shop, userID dot.ID, r *apishop.RequestExportRequest) (_ *apishop.RequestExportResponse, _err error) {
	if userID == 0 {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "")
	}

	// idempotency
	key := shop.ID.String()
	if err := idempgroup.Acquire(key, claim.Token); err != nil {
		return nil, idemp.WrapError(err, "xuất dữ liệu")
	}
	defer func() {
		e := recover()
		if e != nil {
			_err = cm.Errorf(cm.RuntimePanic, nil, "").
				WithMetap("cause", e)
		}
		// release key when error, keep key if export
		if _err != nil {
			idempgroup.ReleaseKey(key, claim.Token)
		}
	}()

	if r.ExportType != PathShopOrders && r.ExportType != PathShopFulfillments {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "export type is not supported")
	}

	var delimiter rune
	switch r.Delimiter {
	case "":
		delimiter = ','
	case ",", ";", "\t":
		delimiter = rune(r.Delimiter[0])
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "invalid delimiter")
	}
	exportOpts := ExportOption{
		Delimiter: delimiter,
		ExcelMode: r.ExcelCompatibleMode,
	}

	from, to, err := cm.ParseDateFromTo(r.DateFrom, r.DateTo)
	if err != nil {
		return nil, err
	}

	tableNameExport := ""

	switch r.ExportType {
	case PathShopOrders:
		tableNameExport = "orders"
	case PathShopFulfillments:
		tableNameExport = "fulfillments"
	}

	exportID := cm.NewBase54ID()
	fileName := fmt.Sprintf(
		"%v_%s_%v_%v_%v", shop.Code, tableNameExport,
		FormatDateShort(from), FormatDateShort(to),
		FormatDateTimeShort(time.Now()))
	midPath := filepath.Join(exportID[:3], exportID)
	zipFileName := fileName + ".zip"
	fullFileName := filepath.Join(midPath, zipFileName)

	exportItem := &model.ExportAttempt{
		ID:           exportID,
		FileName:     zipFileName,
		StoredFile:   fullFileName,
		ExportType:   r.ExportType,
		DownloadURL:  "",
		AccountID:    shop.ID,
		UserID:       userID,
		CreatedAt:    time.Now(),
		RequestQuery: jsonx.MustMarshalToString(r),
		MimeType:     "text/csv",
		Status:       status4.Z,
	}
	resp := &apishop.RequestExportResponse{
		Id:         exportID,
		Filename:   zipFileName,
		ExportType: r.ExportType,
		Status:     status4.Z,
	}

	if err := sqlstore.ExportAttempt(ctx).Create(exportItem); err != nil {
		return nil, err
	}

	switch r.ExportType {
	case PathShopFulfillments:
		// prepare fulfillments for exporting
		query := &shipping.GetFulfillmentExtendedsQuery{
			IDs:          r.Ids,
			ShopIDs:      []dot.ID{shop.ID},
			DateFrom:     from,
			DateTo:       to,
			Filters:      cmapi.ToFilters(r.Filters),
			ResultAsRows: true,
		}

		defer func() {
			// always close the connection when encounting error
			if _err != nil && query.Result.Rows != nil {
				_ = query.Result.Rows.Close()
			}
		}()
		if err := bus.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		if query.Result.Total == 0 {
			return nil, cm.Errorf(cm.ResourceExhausted, nil, "Không có dữ liệu để xuất. Vui lòng thử lại với điều kiện tìm kiếm khác.")
		}

		go ignoreError(exportAndReportProgress(
			func() { idempgroup.ReleaseKey(key, claim.Token) },
			exportItem, fileName, exportOpts,
			query.Result.Total, query.Result.Rows, query.Result.Opts,
			ExportFulfillments,
		))
	case PathShopOrders:
		query := &ordering.GetOrderExtendedsQuery{
			IDs:          r.Ids,
			ShopIDs:      []dot.ID{shop.ID},
			DateFrom:     from,
			DateTo:       to,
			Filters:      cmapi.ToFilters(r.Filters),
			ResultAsRows: true,
		}

		defer func() {
			// always close the connection when encounting error
			if _err != nil && query.Result.Rows != nil {
				_ = query.Result.Rows.Close()
			}
		}()
		if err := bus.Dispatch(ctx, query); err != nil {
			return nil, err
		}
		if query.Result.Total == 0 {
			return nil, cm.Errorf(cm.ResourceExhausted, nil, "Không có dữ liệu để xuất. Vui lòng thử lại với điều kiện tìm kiếm khác.")
		}

		go ignoreError(exportAndReportProgress(
			func() { idempgroup.ReleaseKey(key, claim.Token) },
			exportItem, fileName, exportOpts,
			query.Result.Total, query.Result.Rows, query.Result.Opts,
			ExportOrders,
		))
	}
	return resp, nil
}

func (s *Service) GetExports(ctx context.Context, shopID dot.ID, r *apishop.GetExportsRequest) (*apishop.GetExportsResponse, error) {
	exportAttempts, err := sqlstore.ExportAttempt(ctx).AccountID(shopID).NotYetExpired().List()
	return &apishop.GetExportsResponse{
		ExportItems: convertpb.PbExportAttempts(exportAttempts),
	}, err
}

func ignoreError(err error) {}
