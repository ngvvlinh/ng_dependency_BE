package imcsv

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"runtime/debug"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"o.o/api/main/location"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/ghn_note_code"
	"o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/top/types/etc/status4"
	"o.o/api/top/types/etc/try_on"
	com "o.o/backend/com/main"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	identitymodel "o.o/backend/com/main/identity/model"
	ordermodel "o.o/backend/com/main/ordering/model"
	"o.o/backend/com/main/ordering/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etc/typeutil"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/auth"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
	"o.o/backend/tools/pkg/acl" // TODO: remove this
	"o.o/capi/dot"
	"o.o/common/strs"
	"o.o/common/xerrors"
)

type Import struct {
	uploader         *upload.Uploader
	locationBus      location.QueryBus
	shopVariantStore catalogsqlstore.ShopVariantStoreFactory
}

func New(_locationBus location.QueryBus, rd redis.Store, ul *upload.Uploader, db com.MainDB) (*Import, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60) // 5 minutes
	im := &Import{
		uploader:    ul,
		locationBus: _locationBus,
	}
	im.shopVariantStore = catalogsqlstore.NewShopVariantStore(db)
	if ul != nil {
		ul.ExpectDir(model.ImportTypeShopOrder.String())
	}
	return im, idempgroup.Shutdown
}

func (im *Import) HandleImportOrders(c *httpx.Context) error {
	claim, shop, user, roles := c.SS.Claim(), c.SS.Shop(), c.SS.User(), c.SS.Permission().Roles

	// Do not check permission for 3rd party requests
	if !auth.New().Check(roles, string(acl.ShopOrderImport), shop.IsTest) {
		return cm.Error(cm.PermissionDenied, "", nil)
	}
	key := shop.ID.String()

	resp, _, err := idempgroup.DoAndWrapWithSubkey(c.Context(), key, claim.Token, 30*time.Second, func() (interface{}, error) {
		return im.handleImportOrder(c.Req.Context(), c, shop, user.ID)
	}, "import đơn hàng")
	if err != nil {
		return err
	}

	respMsg := resp.(*types.ImportOrdersResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, claim.Token)
	}
	c.SetResult(respMsg)
	return nil
}

func (im *Import) handleImportOrder(ctx context.Context, c *httpx.Context, shop *identitymodel.Shop, userID dot.ID) (_resp *types.ImportOrdersResponse, _err error) {
	var debugOpts Debug
	if cmenv.NotProd() {
		var err error
		debugOpts, err = parseDebugHeader(c.Req.Header)
		if err != nil {
			return nil, err
		}
	}

	startAt := time.Now()
	imp, err := parseRequest(c)
	if err != nil {
		return nil, err
	}

	file, err := imp.File.Open()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "can not open file")
	}
	defer file.Close()

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "can not open file")
	}

	// We only store file if the file is valid.
	importID := cm.NewIDWithTag(account_tag.TagImport)
	uploadCmd, err := uploadFile(im.uploader, importID, rawData)
	if err != nil {
		return nil, err
	}
	defer func() {
		rerr := recover()
		duration := time.Since(startAt)
		attempt := &model.ImportAttempt{
			ID:           importID,
			UserID:       userID,
			AccountID:    shop.ID,
			OriginalFile: imp.File.Filename,
			StoredFile:   uploadCmd.FileName,
			Type:         model.ImportType(uploadCmd.UploadType),
			NCreated:     0,
			NUpdated:     0,
			NError:       0,
			Status:       0,
			ErrorType:    "",
			Errors:       nil,
			DurationMs:   int(duration / time.Millisecond),
			CreatedAt:    startAt,
		}
		switch {
		case rerr != nil:
			stack := debug.Stack()
			attempt.Status = status4.N
			attempt.ErrorType = "panic"
			savedErr := cm.Errorf(cm.Internal, nil, "%v", rerr).
				WithMeta("stack", cm.UnsafeBytesToString(stack))
			attempt.Errors = []*model.Error{model.ToError(savedErr)}

			// respond internal error to client
			_err = cm.Error(cm.Internal, "", nil)

		case _err != nil:
			attempt.Status = status4.N
			attempt.ErrorType = "error"
			_err = xerrors.ToError(_err).WithMetaID("import_id", importID)
			attempt.Errors = []*model.Error{model.ToError(_err)}

		case len(_resp.CellErrors) > 0:
			attempt.Status = status4.N
			attempt.ErrorType = "cell_errors"
			attempt.Errors = cmapi.ErrorsToModel(_resp.CellErrors)
			attempt.NError = len(_resp.CellErrors)

		case len(_resp.ImportErrors) > 0:
			count := cmapi.CountErrors(_resp.ImportErrors)
			if count == 0 {
				attempt.Status = status4.P
				attempt.NCreated = len(_resp.ImportErrors)

			} else {
				attempt.Status = status4.S // partially error
				attempt.ErrorType = "import_errors"
				attempt.Errors = cmapi.ErrorsToModel(_resp.ImportErrors)
				attempt.NError = count
				attempt.NCreated = len(_resp.ImportErrors) - count
			}

		default:
			attempt.Status = 0 // unknown
		}
		if _resp != nil {
			_resp.ImportId = importID
		}

		createAttemptCmd := &model.CreateImportAttemptCommand{
			ImportAttempt: attempt,
		}
		if err = bus.Dispatch(ctx, createAttemptCmd); err != nil {
			_err = err
		}
	}()

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "invalid file format")
	}

	sheetName, err := validateSheets(excelFile)
	if err != nil {
		return nil, err
	}
	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}
	cleanRows(rows)
	imp.Rows = rows

	idx, _errs, err := schema.ValidateSchema(ctx, &rows[0])
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}
	imp.LastRow, _errs, err = validateRows(ctx, idx, rows, idxOrderEdCode, idxLines)
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}

	rowOrders, _errs, err := parseRows(idx, imp.Mode, imp.GHNNoteCode, rows)
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}
	now := time.Now()
	orders := make([]*ordermodel.Order, len(rowOrders))
	for i, rowOrder := range rowOrders {
		order, errs := parseRowToModel(im.locationBus, idx, imp.Mode, shop, rowOrder, now, userID)
		if len(errs) > 0 {
			_errs = append(_errs, errs...)
			if len(_errs) >= MaxCellErrors {
				break
			}
			continue
		}
		orders[i] = order
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}

	_errs, err = im.verifyOrders(ctx, shop, idx, imp.CodeMode, rowOrders)
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}

	// Fill product and variant ids
	// Because we call parseRowToModel() before VerifyOrders()
	for i, rowOrder := range rowOrders {
		order := orders[i]
		order.IsOutsideEtop = false

		var productIDs, variantIDs []dot.ID
		for j, rowOrderLine := range rowOrder.Lines {
			line := order.Lines[j]

			if rowOrderLine.XVariant != nil {
				line.ProductName = rowOrderLine.XVariant.GetFullName()
				line.VariantID = rowOrderLine.XVariant.ShopVariant.VariantID
				line.ProductID = rowOrderLine.XVariant.ShopProduct.ProductID
				productIDs = append(productIDs, line.ProductID)
				variantIDs = append(variantIDs, line.VariantID)
			} else {
				order.IsOutsideEtop = true
			}
		}
		order.ProductIDs = productIDs
		order.VariantIDs = variantIDs
	}

	_errs = make([]error, len(orders))
	for i, order := range orders {
		if debugOpts.FailPercent != 0 && isRandomFail(debugOpts.FailPercent) {
			_errs[i] = cm.NSErrorf(cm.Internal, nil, "Random error for development")
			continue
		}

		cmd := &modelx.CreateOrderCommand{
			Order: order,
		}
		err := bus.Dispatch(ctx, cmd)
		_errs[i] = err
	}
	if xerrors.Errors(_errs).IsAll() {
		return nil, cm.Errorf(cm.Internal, _errs[0], "Không thể import đơn hàng. Vui lòng liên hệ %v.", wl.X(ctx).CSEmail)
	}

	resp := &types.ImportOrdersResponse{
		Data:         imp.toSpreadsheetData(idx),
		Orders:       convertpb.PbOrders(orders, account_tag.TagShop),
		ImportErrors: cmapi.PbErrors(_errs),
	}
	// Remove failed order from the response
	for i, err := range _errs {
		if err != nil {
			resp.Orders[i] = nil
		}
	}
	return resp, nil
}

func validateSheets(file *excelize.File) (sheetName string, err error) {
	sheetName = file.GetSheetName(1)
	if sheetName == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Không thể đọc được file.").WithMeta("reason", "invalid sheet")
	}

	norm := validate.NormalizeSearchSimple(sheetName)
	if !strings.Contains(norm, "don hang") {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Sheet đầu tiên trong file phải là danh sách đơn hàng cần import.").WithMeta("reason", "invalid sheet name")
	}

	return sheetName, nil
}

func validateRows(ctx context.Context, idx imcsv.Indexer, rows [][]string, idxCode, idxLines int) (lastNonEmptyRow int, errs []error, _ error) {
	if len(rows) > MaxRows {
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "File import quá lớn. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	var orderRow []string
	for r := 1; r < len(rows); r++ {
		if len(rows[r]) < len(schema) {
			err := imcsv.CellError(idx, r, -1, "Số cột không đúng cấu trúc yêu cầu.")
			errs = append(errs, err)
			if len(errs) >= MaxCellErrors {
				return
			}
			continue
		}

		imcsv.CleanRow(&rows[r], len(schema))

		preRow := rows[r-1]
		preCode := idx.GetCell(preRow, idxCode)

		row := rows[r]
		code := idx.GetCell(row, idxCode)

		if code != "" {
			if r-lastNonEmptyRow > MaxEmptyRows {
				return lastNonEmptyRow, nil, imcsv.CellError(idx, r, -1, "Dòng %v: Số dòng để trống liên tiếp tối đa là %v", MaxEmptyRows)
			}
			lastNonEmptyRow = r
		}

		switch {
		case code == "":
			// All cells must be empty
			for c := 0; c < len(row); c++ {
				if c != idxUnderscore && idx.GetCell(row, c) != "" {
					err := imcsv.CellError(idx, r, c, "Giá trị phải để trống (vì mã đơn hàng để trống)")
					errs = append(errs, err)
					if len(errs) >= MaxCellErrors {
						return
					}
					continue
				}
			}

		case code != preCode:
			orderRow = row

			nLines, err := imcsv.ParseUint(idx.GetCell(row, idxLines))
			if err != nil {
				err = imcsv.CellError(idx, r, idxLines, err.Error())
				errs = append(errs, err)
				if len(errs) >= MaxCellErrors {
					return
				}
				continue // skip the following check
			}
			if err := checkOrderRowNLines(idx, r, rows, idxCode, idxLines, nLines); err != nil {
				errs = append(errs, err)
				if len(errs) >= MaxCellErrors {
					return
				}
			}

		case code == preCode:
			if orderRow == nil {
				return lastNonEmptyRow, nil, imcsv.CellError(idx, r, idxCode, "Mã đơn hàng không hợp lệ. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail)
			}

			if idx.GetCell(row, idxLines) != "" {
				err := imcsv.CellError(idx, r, idxLines, "Giá trị phải để trống (vì thuộc cùng một đơn hàng)")
				errs = append(errs, err)
				if len(errs) >= MaxCellErrors {
					return
				}
			}

			// All cells (not line and hidden) must be empty or equal to the previous line
			for c, col := range schema {
				val := idx.GetCell(row, c)
				if !col.Line && !col.Hidden && val != "" && val != idx.GetCell(orderRow, c) {
					err := imcsv.CellError(idx, r, c, "Giá trị phải để trống hoặc bằng giá trị ở dòng trên (vì thuộc cùng một đơn hàng)")
					errs = append(errs, err)
					if len(errs) >= MaxCellErrors {
						return
					}
				}
			}
		}
	}

	if orderRow == nil {
		return lastNonEmptyRow, nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}
	return
}

func checkOrderRowNLines(idx imcsv.Indexer, r int, rows [][]string, idxCode, idxLines, nLines int) error {
	if nLines <= 0 {
		return imcsv.CellError(idx, r, idxLines, "Số dòng sản phẩm không hợp lệ.")
	}
	if r+nLines > len(rows) {
		return imcsv.CellError(idx, r, idxLines, "Số dòng sản phẩm không đúng.").
			WithMeta("reason", "overflow")
	}

	orderCode := rows[r][idxCode]
	for i := r + 1; i < r+nLines-1; i++ {
		code := rows[r][idxCode]
		if code != orderCode {
			return imcsv.CellError(idx, r, idxLines, "Số dòng sản phẩm không đúng").
				WithMeta("reason", "not match")
		}
	}
	return nil
}

func parseRows(idx imcsv.Indexer, mode Mode, ghnNoteCode ghn_note_code.GHNNoteCode, rows [][]string) (orders []*RowOrder, _errs []error, _ error) {
	var currentRowOrder, rowOrder *RowOrder

	// Skip the first row
	for r := 1; r < len(rows); r++ {
		row := rows[r]
		rowOrder, _errs = parseRow(idx, mode, ghnNoteCode, r, row, currentRowOrder, _errs)
		if rowOrder != currentRowOrder && rowOrder != nil {
			orders = append(orders, rowOrder)
			currentRowOrder = rowOrder
		}
	}
	return orders, _errs, nil
}

func parseRow(idx imcsv.Indexer, mode Mode, ghnNoteCode ghn_note_code.GHNNoteCode, r int, row []string, currentRowOrder *RowOrder, _errs []error) (*RowOrder, []error) {
	code := idx.GetCell(row, idxOrderEdCode)
	if code == "" {
		return nil, _errs
	}

	nLines := idx.GetCell(row, idxLines)
	switch {
	case nLines != "":
		newRowOrder, errs := parseRowAsOrder(idx, mode, ghnNoteCode, r, row)
		_errs = append(_errs, errs...)
		if len(_errs) >= MaxCellErrors {
			return newRowOrder, _errs
		}

		line, errs := parseRowAsOrderLine(idx, mode, r, row)
		newRowOrder.Lines = append(newRowOrder.Lines, line)
		_errs = append(_errs, errs...)
		return newRowOrder, _errs

	default:
		if currentRowOrder == nil {
			err := imcsv.CellError(idx, r, -1, "Unexpected: Must be the first line of an order")
			_errs = append(_errs, err)
			return nil, _errs
		}

		line, errs := parseRowAsOrderLine(idx, mode, r, row)
		currentRowOrder.Lines = append(currentRowOrder.Lines, line)
		_errs = append(_errs, errs...)
		return currentRowOrder, _errs
	}
}

func parseRowAsOrder(idx imcsv.Indexer, mode Mode, ghnNoteCode ghn_note_code.GHNNoteCode, r int, row []string) (_ *RowOrder, errs []error) {
	var col int
	var rowOrder RowOrder

	col = idxTotalItems
	totalItems, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxTotalWeight
	totalWeight, err := imcsv.ParseFloat(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	totalWeightInt := int(math.Floor(totalWeight * 1000))
	if totalWeightInt%10 == 9 {
		totalWeightInt++
	}

	col = idxShippingNoteGhn
	if ghnNoteCode == 0 {
		ghnNoteCode, err = parseAsGHNNoteCode(idx.GetCell(row, col))
		if err != nil {
			err = imcsv.CellError(idx, r, col, err.Error())
			errs = append(errs, err)
		}
	}

	col = idxBasketValue
	basketValue, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxBasketValueDiscounted
	basketValueDiscounted, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxOrderDiscount
	orderDiscount, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxFeeLineShipping
	rowOrder.FeeLineShipping, err = imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxFeeLineTax
	rowOrder.FeeLineTax, err = imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxFeeLineOther
	rowOrder.FeeLineOther, err = imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxTotalFee
	rowOrder.TotalFee, err = imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}
	expectedTotalFee := rowOrder.FeeLineShipping + rowOrder.FeeLineTax + rowOrder.FeeLineOther
	if err == nil && rowOrder.TotalFee != expectedTotalFee {
		err = imcsv.CellError(idx, r, col, "Tổng phí không đúng (cần là tổng của các cột chi phí, mong đợi: %v, giá trị đã nhập: %v)", expectedTotalFee, rowOrder.TotalFee)
		errs = append(errs, err)
	}

	col = idxTotalAmount
	totalAmount, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}
	expectedTotalAmount := basketValueDiscounted - orderDiscount + expectedTotalFee
	if err == nil && totalAmount != expectedTotalAmount {
		err = imcsv.CellError(idx, r, col, "Tổng tiền thanh toán không đúng (tổng tiền hàng sau giảm giá: %v, giảm giá đơn hàng: %v, tổng phí: %v, mong đợi: %v, đã nhập: %v)", basketValueDiscounted, orderDiscount, expectedTotalFee, expectedTotalAmount, totalAmount)
		errs = append(errs, err)
	}

	col = idxShopCod
	shopCOD, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxIsCod
	isCod, err := imcsv.ParseBool(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}
	if isCod != (shopCOD != 0) {
		err = imcsv.CellError(idx, r, col, "Thu hộ không hợp lệ. Giá trị tiền thu hộ là %v nhưng \"Thu hộ\" là \"%v\"", shopCOD, idx.GetCell(row, col))
		errs = append(errs, err)
	}

	rowOrder.RowIndex = r
	rowOrder.OrderEdCode = idx.GetCell(row, idxOrderEdCode)
	rowOrder.CustomerName = idx.GetCell(row, idxCustomerName)
	rowOrder.CustomerPhone = idx.GetCell(row, idxCustomerPhone)
	rowOrder.ShippingAddress = idx.GetCell(row, idxShippingAddress)
	rowOrder.ShippingProvince = idx.GetCell(row, idxShippingProvince)
	rowOrder.ShippingDistrict = idx.GetCell(row, idxShippingDistrict)
	rowOrder.ShippingWard = idx.GetCell(row, idxShippingWard)
	rowOrder.ShippingNote = idx.GetCell(row, idxShippingNote)
	rowOrder.GHNNoteCode = ghnNoteCode
	rowOrder.BasketValue = basketValue
	rowOrder.BasketValueDiscounted = basketValueDiscounted
	rowOrder.OrderDiscount = orderDiscount
	rowOrder.TotalAmount = totalAmount
	rowOrder.ShopCOD = shopCOD
	rowOrder.IsCOD = isCod
	rowOrder.TotalWeight = totalWeightInt
	rowOrder.TotalItems = totalItems
	rowOrder.Lines = nil
	return &rowOrder, errs
}

func parseRowAsOrderLine(idx imcsv.Indexer, mode Mode, r int, row []string) (_ *RowOrderLine, errs []error) {
	var col int

	col = idxLineQuantity
	quantity, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	} else if quantity == 0 {
		err = imcsv.CellError(idx, r, col, "Số lượng bằng 0")
		errs = append(errs, err)
	}

	col = idxVariantRetailPrice
	retailPrice, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)

	} else if retailPrice == 0 {
		err = imcsv.CellError(idx, r, col, "Đơn giá bằng 0")
		errs = append(errs, err)
	}

	col = idxLineAmount
	lineAmount, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)

	} else if retailPrice*quantity != lineAmount {
		err = imcsv.CellError(idx, r, col, "Thành tiền (trước giảm giá) không đúng (đơn giá: %v, số lượng: %v, thành tiền: %v)", retailPrice, quantity, lineAmount)
		errs = append(errs, err)
	}

	col = idxLineDiscountPercent
	var lineDiscountPercent float64
	if val := idx.GetCell(row, col); val != "" {
		lineDiscountPercent, err = imcsv.ParsePercent(val)
		if err != nil {
			err = imcsv.CellError(idx, r, col, err.Error())
			errs = append(errs, err)
		}
	}

	col = idxLineDiscountValue
	var lineDiscountValue int
	if val := idx.GetCell(row, col); val != "" {
		lineDiscountValue, err = imcsv.ParseUint(val)
		if err != nil {
			err = imcsv.CellError(idx, r, col, err.Error())
			errs = append(errs, err)
		}
	}

	col = idxLineTotalAmount
	lineTotalAmount, err := imcsv.ParseUint(idx.GetCell(row, col))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)

	} else if lineTotalAmount != lineAmount-
		int(float64(lineAmount)*lineDiscountPercent)-
		lineDiscountValue*quantity {
		err = imcsv.CellError(idx, r, col, "Thành tiền (sau giảm giá) không đúng (trước giảm giá: %v, giảm giá %%: %v, giảm giá đ: %v, sau giảm giá: %v)", lineAmount, lineDiscountPercent*100, lineDiscountValue, lineTotalAmount)
		errs = append(errs, err)
	}

	line := &RowOrderLine{
		RowIndex: r,

		VariantEdCode:       idx.GetCell(row, idxVariantEdCode),
		VariantName:         idx.GetCell(row, idxVariantName),
		RetailPrice:         retailPrice,
		Quantity:            quantity,
		LineAmount:          lineAmount,
		LineDiscountPercent: lineDiscountPercent,
		LineDiscountValue:   lineDiscountValue,
		LineTotalAmount:     lineTotalAmount,
	}
	return line, errs
}

func parseAsGHNNoteCode(v string) (ghn_note_code.GHNNoteCode, error) {
	switch v {
	case "":
		return 0, nil
	case "Cho xem hàng không thử",
		ghn_note_code.CHOXEMHANGKHONGTHU.String():
		return ghn_note_code.CHOXEMHANGKHONGTHU, nil
	case "Cho thử hàng",
		ghn_note_code.CHOTHUHANG.String():
		return ghn_note_code.CHOTHUHANG, nil
	case "Không cho xem hàng",
		ghn_note_code.KHONGCHOXEMHANG.String():
		return ghn_note_code.KHONGCHOXEMHANG, nil
	}

	ghnNote := validate.NormalizeSearchSimple(v)
	switch ghnNote {
	case "":
		return 0, nil
	case try_on.NoteChoXemHang:
		return ghn_note_code.CHOXEMHANGKHONGTHU, nil
	case try_on.NoteChoThuHang:
		return ghn_note_code.CHOTHUHANG, nil
	case try_on.NoteKhongXemHang:
		return ghn_note_code.KHONGCHOXEMHANG, nil
	default:
		return 0, errors.New("Ghi chú xem hàng không hợp lệ, cần một trong các giá trị: 'Cho thử hàng', 'Cho xem hàng không thử', 'Không cho xem hàng'.")
	}
}

func parseRowToModel(locationBus location.QueryBus, idx imcsv.Indexer, mode Mode, shop *identitymodel.Shop, rowOrder *RowOrder, now time.Time, userID dot.ID) (_ *ordermodel.Order, _errs []error) {
	_errs = rowOrder.Validate(idx, mode)
	address, err := parseAddress(locationBus, rowOrder)
	if err != nil {
		err = imcsv.CellError(idx, rowOrder.RowIndex+1, -1, err.Error())
		_errs = append(_errs, err)
	}

	totalItems, basketValue, totalLineDiscount := 0, 0, 0
	lines := make([]*ordermodel.OrderLine, len(rowOrder.Lines))
	for i, rowOrderLine := range rowOrder.Lines {
		line, errs := parseLineToModel(idx, mode, rowOrderLine)
		if len(errs) > 0 {
			_errs = append(_errs, errs...)
			continue
		}
		line.ShopID = shop.ID

		lines[i] = line
		totalItems += line.Quantity
		basketValue += line.RetailPrice * line.Quantity
		totalLineDiscount += line.TotalDiscount
	}
	if len(_errs) > 0 {
		return nil, _errs
	}

	order := &ordermodel.Order{
		CreatedBy:                 userID,
		ID:                        0, // will be filled by sqlstore
		ShopID:                    shop.ID,
		Code:                      "",
		EdCode:                    rowOrder.OrderEdCode,
		ProductIDs:                nil, // will be filled later
		VariantIDs:                nil, // will be filled later
		Currency:                  model.CurrencyVND,
		PaymentMethod:             0,
		Customer:                  parseCustomer(rowOrder),
		CustomerAddress:           address,
		BillingAddress:            address,
		ShippingAddress:           address,
		CustomerPhone:             rowOrder.CustomerPhone,
		CustomerEmail:             "",
		CreatedAt:                 now,
		ProcessedAt:               now,
		UpdatedAt:                 time.Time{},
		ClosedAt:                  time.Time{},
		ConfirmedAt:               time.Time{},
		CancelledAt:               time.Time{},
		CancelReason:              "",
		CustomerConfirm:           0,
		ShopConfirm:               0,
		ConfirmStatus:             0,
		FulfillmentShippingStatus: 0,
		Status:                    0,
		Lines:                     lines,
		FeeLines:                  rowOrder.GetFeeLines(),
		TotalFee:                  rowOrder.TotalFee,
		Discounts:                 nil,
		TotalItems:                totalItems,
		BasketValue:               basketValue,
		TotalWeight:               rowOrder.TotalWeight,
		TotalTax:                  0,
		OrderDiscount:             rowOrder.OrderDiscount,
		TotalDiscount:             rowOrder.OrderDiscount + totalLineDiscount,
		ShopShippingFee:           rowOrder.FeeLineShipping,
		ShopCOD:                   rowOrder.ShopCOD,
		TotalAmount:               rowOrder.TotalAmount,
		OrderNote:                 "",
		ShopNote:                  "",
		ShippingNote:              rowOrder.ShippingNote,
		OrderSourceType:           order_source.Import,
		OrderSourceID:             0,
		ExternalOrderID:           "",
		ReferenceURL:              "",
		ShopShipping:              nil,
		IsOutsideEtop:             false, // will be filled later
		GhnNoteCode:               rowOrder.GHNNoteCode,
		TryOn:                     typeutil.TryOnFromGHNNoteCode(rowOrder.GHNNoteCode),
	}
	if rowOrder.ShopCOD != 0 {
		order.PaymentMethod = payment_method.COD
	} else {
		order.PaymentMethod = payment_method.Other
	}
	return order, _errs
}

func parseLineToModel(idx imcsv.Indexer, mode Mode, rowOrderLine *RowOrderLine) (*ordermodel.OrderLine, []error) {
	errs := rowOrderLine.Validate(idx, mode)
	if len(errs) > 0 {
		return nil, errs
	}

	line := &ordermodel.OrderLine{
		OrderID:         0, // will be filled when insert
		VariantID:       0,
		ProductName:     rowOrderLine.VariantName,
		ProductID:       0,
		ShopID:          0, // will be filled later
		Weight:          0,
		Quantity:        rowOrderLine.Quantity,
		ListPrice:       0,
		RetailPrice:     rowOrderLine.RetailPrice,
		PaymentPrice:    rowOrderLine.PaymentPrice,
		LineAmount:      rowOrderLine.LineAmount,
		TotalDiscount:   rowOrderLine.LineAmount - rowOrderLine.LineTotalAmount,
		TotalLineAmount: rowOrderLine.LineTotalAmount,
		ImageURL:        "",
		Attributes:      nil,
		IsOutsideEtop:   false,
		Code:            rowOrderLine.VariantEdCode,
	}
	if rowOrderLine.XVariant != nil {
		line.ProductID = rowOrderLine.XVariant.ShopProduct.ProductID
		line.VariantID = rowOrderLine.XVariant.ShopVariant.VariantID
	}
	return line, nil
}

func parseCustomer(rowOrder *RowOrder) *ordermodel.OrderCustomer {
	return &ordermodel.OrderCustomer{
		FirstName:     "",
		LastName:      "",
		FullName:      rowOrder.CustomerName,
		Email:         "",
		Gender:        "",
		Birthday:      "",
		VerifiedEmail: false,
		ExternalID:    "",
	}
}

func parseAddress(locationBus location.QueryBus, rowOrder *RowOrder) (*ordermodel.OrderAddress, error) {
	var loc *location.LocationQueryResult
	var err error
	if rowOrder.ShippingProvince != "" || rowOrder.ShippingDistrict != "" || rowOrder.ShippingWard != "" {
		loc, err = parseLocation(locationBus, rowOrder)
		if err != nil {
			return nil, err
		}
	}

	address := &ordermodel.OrderAddress{
		FullName:     rowOrder.CustomerName,
		FirstName:    "",
		LastName:     "",
		Phone:        rowOrder.CustomerPhone,
		Country:      location.CountryVietnam,
		City:         "",
		Province:     "",
		District:     "",
		Ward:         "",
		Zip:          "",
		DistrictCode: "",
		ProvinceCode: "",
		WardCode:     "",
		Company:      "",
		Address1:     rowOrder.ShippingAddress,
		Address2:     "",
	}
	if loc == nil {
		return address, nil
	}
	if loc.Province != nil {
		address.Province = loc.Province.Name
		address.ProvinceCode = loc.Province.Code
	}
	if loc.District != nil {
		address.District = loc.District.Name
		address.DistrictCode = loc.District.Code
	}
	if loc.Ward != nil {
		address.Ward = loc.Ward.Name
		address.WardCode = loc.Ward.Code
	} else {
		address.Ward = rowOrder.ShippingWard
	}
	return address, nil
}

func parseLocation(locationBus location.QueryBus, rowOrder *RowOrder) (*location.LocationQueryResult, error) {
	query := &location.FindLocationQuery{
		Province: rowOrder.ShippingProvince,
		District: rowOrder.ShippingDistrict,
		Ward:     rowOrder.ShippingWard,
	}
	if err := locationBus.Dispatch(context.TODO(), query); err != nil {
		return nil, err
	}
	loc := query.Result

	if loc.Province == nil {
		return nil, fmt.Errorf(
			"Địa chỉ không hợp lệ: %v, %v, %v (không tìm thấy tỉnh/thành phố).",
			strs.TrimMax(rowOrder.ShippingWard, 200),
			strs.TrimMax(rowOrder.ShippingDistrict, 200),
			strs.TrimMax(rowOrder.ShippingProvince, 200),
		)
	}
	if loc.District == nil {
		return nil, fmt.Errorf(
			"Địa chỉ không hợp lệ: %v, %v, %v (không tìm thấy quận/huyện).",
			strs.TrimMax(rowOrder.ShippingWard, 200),
			strs.TrimMax(rowOrder.ShippingDistrict, 200),
			strs.TrimMax(rowOrder.ShippingProvince, 200),
		)
	}
	if loc.Ward == nil {
		return nil, fmt.Errorf(
			"Địa chỉ không hợp lệ: %v, %v, %v (không tìm thấy phường/xã).",
			strs.TrimMax(rowOrder.ShippingWard, 200),
			strs.TrimMax(rowOrder.ShippingDistrict, 200),
			strs.TrimMax(rowOrder.ShippingProvince, 200),
		)
	}
	return loc, nil
}
