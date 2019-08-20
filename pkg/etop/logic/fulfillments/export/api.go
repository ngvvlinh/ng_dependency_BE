package export

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"etop.vn/backend/com/main/shipping/modelx"
	pbcm "etop.vn/backend/pb/common"
	pbs4 "etop.vn/backend/pb/etop/etc/status4"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	wrappershop "etop.vn/backend/wrapper/etop/shop"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("export",
		RequestExport,
		GetExports,
	)
}

func RequestExport(ctx context.Context, r *wrappershop.RequestExportEndpoint) (_err error) {
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

	if r.ExportType != PathShopFulfillments {
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

	// prepare fulfillments for exporting
	query := &modelx.GetFulfillmentExtendedsQuery{
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

	exportID := cm.NewBase54ID()
	fileName := fmt.Sprintf(
		"%v_fulfillments_%v_%v_%v", shop.Code,
		formatDateShort(from), formatDateShort(to),
		formatDateTimeShort(time.Now()))
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

	go exportFulfillmentsAndReportProgress(
		func() { idempgroup.ReleaseKey(key1, claim.Token) },
		exportItem, fileName, exportOpts,
		query.Result.Total, query.Result.Rows, query.Result.Opts,
	)

	r.Result = resp
	return nil
}

func GetExports(ctx context.Context, r *wrappershop.GetExportsEndpoint) error {
	shopID := r.Context.Shop.ID
	exportAttempts, err := sqlstore.ExportAttempt(ctx).AccountID(shopID).NotYetExpired().List()

	r.Result = &pbshop.GetExportsResponse{
		ExportItems: pbshop.PbExportAttempts(exportAttempts),
	}
	return err
}
