package export

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	ordering "etop.vn/backend/com/main/ordering/modelx"
	shipping "etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrappershop "etop.vn/backend/wrapper/etop/shop"
)

func init() {
	bus.AddHandlers("export",
		s.RequestExport,
		s.GetExports,
	)
}

var s = &Service{}

type Service struct{}

func (s *Service) RequestExport(ctx context.Context, r *wrappershop.RequestExportEndpoint) (_err error) {
	claim := r.Context.Claim
	shop := r.Context.Shop
	userID := r.Context.UserID
	if userID == 0 {
		return cm.Errorf(cm.PermissionDenied, nil, "")
	}

	// idempotency
	key1 := strconv.FormatInt(shop.ID, 10)
	if err := idempgroup.Acquire(key1, claim.Token); err != nil {
		return idemp.WrapError(err, "xuất dữ liệu")
	}
	defer func() {
		e := recover()
		if e != nil {
			_err = cm.Errorf(cm.RuntimePanic, nil, "").
				WithMetap("cause", e)
		}
		// release key when error, keep key if export
		if _err != nil {
			idempgroup.ReleaseKey(key1, claim.Token)
		}
	}()

	if r.ExportType != PathShopOrders && r.ExportType != PathShopFulfillments {
		return cm.Errorf(cm.InvalidArgument, nil, "export type is not supported")
	}

	var delimiter rune
	switch r.Delimiter {
	case "":
		delimiter = ','
	case ",", ";", "\t":
		delimiter = rune(r.Delimiter[0])
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "invalid delimiter")
	}
	exportOpts := ExportOption{
		Delimiter: delimiter,
		ExcelMode: r.ExcelCompatibleMode,
	}

	from, to, err := cm.ParseDateFromTo(r.DateFrom, r.DateTo)
	if err != nil {
		return err
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
		RequestQuery: pbcm.MustMarshalToString(r.RequestExportRequest),
		MimeType:     "text/csv",
		Status:       model.S4Zero,
	}
	resp := &pbshop.RequestExportResponse{
		Id:         exportID,
		Filename:   zipFileName,
		ExportType: r.ExportType,
		Status:     pbs4.Pb(model.S4Zero),
	}

	if err := sqlstore.ExportAttempt(ctx).Create(exportItem); err != nil {
		return err
	}

	switch r.ExportType {
	case PathShopFulfillments:
		// prepare fulfillments for exporting
		query := &shipping.GetFulfillmentExtendedsQuery{
			IDs:          r.Ids,
			ShopIDs:      []int64{shop.ID},
			DateFrom:     from,
			DateTo:       to,
			Filters:      pbcm.ToFilters(r.Filters),
			ResultAsRows: true,
		}

		defer func() {
			// always close the connection when encounting error
			if _err != nil && query.Result.Rows != nil {
				_ = query.Result.Rows.Close()
			}
		}()
		if err := bus.Dispatch(ctx, query); err != nil {
			return err
		}
		if query.Result.Total == 0 {
			return cm.Errorf(cm.ResourceExhausted, nil, "Không có dữ liệu để xuất. Vui lòng thử lại với điều kiện tìm kiếm khác.")
		}

		go ignoreError(exportAndReportProgress(
			func() { idempgroup.ReleaseKey(key1, claim.Token) },
			exportItem, fileName, exportOpts,
			query.Result.Total, query.Result.Rows, query.Result.Opts,
			ExportFulfillments,
		))
	case PathShopOrders:
		query := &ordering.GetOrderExtendedsQuery{
			IDs:          r.Ids,
			ShopIDs:      []int64{shop.ID},
			DateFrom:     from,
			DateTo:       to,
			Filters:      pbcm.ToFilters(r.Filters),
			ResultAsRows: true,
		}

		defer func() {
			// always close the connection when encounting error
			if _err != nil && query.Result.Rows != nil {
				_ = query.Result.Rows.Close()
			}
		}()
		if err := bus.Dispatch(ctx, query); err != nil {
			return err
		}
		if query.Result.Total == 0 {
			return cm.Errorf(cm.ResourceExhausted, nil, "Không có dữ liệu để xuất. Vui lòng thử lại với điều kiện tìm kiếm khác.")
		}

		go ignoreError(exportAndReportProgress(
			func() { idempgroup.ReleaseKey(key1, claim.Token) },
			exportItem, fileName, exportOpts,
			query.Result.Total, query.Result.Rows, query.Result.Opts,
			ExportOrders,
		))
	}

	r.Result = resp
	return nil
}

func (s *Service) GetExports(ctx context.Context, r *wrappershop.GetExportsEndpoint) error {
	shopID := r.Context.Shop.ID
	exportAttempts, err := sqlstore.ExportAttempt(ctx).AccountID(shopID).NotYetExpired().List()

	r.Result = &pbshop.GetExportsResponse{
		ExportItems: pbshop.PbExportAttempts(exportAttempts),
	}
	return err
}

func ignoreError(err error) {}
