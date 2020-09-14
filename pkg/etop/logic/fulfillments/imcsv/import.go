package imcsv

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/status4"
	identitymodel "o.o/backend/com/main/identity/model"
	"o.o/backend/com/main/location"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/etop/upload"
	"o.o/capi/dot"
	"o.o/common/xerrors"
)

type Import struct {
	Uploader *upload.Uploader

	ExportAttemptStore sqlstore.ExportAttemptStoreInterface
}

func New(rd redis.Store, ul *upload.Uploader, exportAttemptStore sqlstore.ExportAttemptStoreInterface) (*Import, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60) // 5 minutes
	im := &Import{
		Uploader:           ul,
		ExportAttemptStore: exportAttemptStore,
	}
	return im, idempgroup.Shutdown
}

func (im *Import) HandleImportFulfillments(c *httpx.Context) error {
	claim, shop, user, _ := c.SS.Claim(), c.SS.Shop(), c.SS.User(), c.SS.Permission().Roles

	// TODO: Check permission
	key := shop.ID.String()

	resp, _, err := idempgroup.DoAndWrapWithSubkey(c.Context(), key, claim.Token, 30*time.Second, func() (interface{}, error) {
		return im.handleImportFulfillments(c.Req.Context(), c, shop, user.ID)
	}, "import đơn giao hàng")
	if err != nil {
		return err
	}
	respMsg := resp.(*types.ImportFulfillmentsResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, claim.Token)
	}
	c.SetResult(respMsg)
	return nil
}

func (im *Import) handleImportFulfillments(ctx context.Context, c *httpx.Context, shop *identitymodel.Shop, userID dot.ID) (_resp *types.ImportFulfillmentsResponse, _err error) {
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
	uploadCmd, err := uploadFile(im.Uploader, importID, rawData)
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

		default:
			attempt.Status = 0 // unknown
		}
		if _resp != nil {
			_resp.ImportID = importID
		}

		createAttemptCmd := &model.CreateImportAttemptCommand{
			ImportAttempt: attempt,
		}
		if err = im.ExportAttemptStore.CreateImportAttempt(ctx, createAttemptCmd); err != nil {
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

	rows = cleanEmptyRows(rows)
	cleanRows(rows)

	imp.Rows = rows

	idx, _errs, err := schema.ValidateSchema(ctx, &rows[0])
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(idx, _errs)
	}

	imp.LastRow, _errs, err = validateRows(ctx, idx, rows)
	if err != nil {
		return nil, err
	}

	rowFulfillments, _errs, err := parseRows(idx, rows)
	if err != nil {
		return nil, err
	}

	_errs, err = im.verifyFulfillments(ctx, shop, idx, rowFulfillments)
	if err != nil {
		return nil, err
	}

	resp := &types.ImportFulfillmentsResponse{
		Data:            imp.toSpreadsheetData(idx),
		SpecificColumns: imp.toSpecificColumns(idx),
		Fulfillments:    parseRowsToModels(rowFulfillments),
		CellErrors:      cmapi.PbErrors(_errs),
	}

	return resp, nil
}

func parseRequest(c *httpx.Context) (*Importer, error) {
	// Limit the max file size
	c.Req.Body = http.MaxBytesReader(c.Resp, c.Req.Body, MaxFilesize)

	form, err := c.MultipartForm()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
		// continue
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	req := &Importer{
		Schema: schema,
		File:   files[0],
	}
	return req, nil
}

func validateSheets(file *excelize.File) (sheetName string, err error) {
	sheetName = file.GetSheetName(1)
	if sheetName == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Không thể đọc được file.").WithMeta("reason", "invalid sheet")
	}

	norm := validate.NormalizeSearchSimple(sheetName)
	if !strings.Contains(norm, "don giao hang") {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Sheet đầu tiên trong file phải là danh sách đơn giao hàng cần import.").WithMeta("reason", "invalid sheet name")
	}

	return sheetName, nil
}

func validateRows(ctx context.Context, idx imcsv.Indexer, rows [][]string) (lastNonEmptyRow int, errs []error, _ error) {
	if len(rows) > MaxRows {
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "File import quá lớn. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail)
	}

	for r := 1; r < len(rows); r++ {
		if len(rows[r]) < len(schema) {
			err := imcsv.CellError(idx, r, -1, "Số cột không đúng cấu trúc yêu cầu.")
			errs = append(errs, err)
			continue
		}

		imcsv.CleanRow(&rows[r], len(schema))

		lastNonEmptyRow = r
	}

	return
}

func parseRows(idx imcsv.Indexer, rows [][]string) (fulfillments []*RowFulfillment, _errs []error, _ error) {
	var currentRowFulfillment, rowFulfillment *RowFulfillment

	// Skip the first row
	for r := 1; r < len(rows); r++ {
		row := rows[r]
		rowFulfillment, _errs = parseRow(idx, r, row, currentRowFulfillment, _errs)
		if rowFulfillment != currentRowFulfillment && rowFulfillment != nil {
			fulfillments = append(fulfillments, rowFulfillment)
			currentRowFulfillment = rowFulfillment
		}
	}
	return fulfillments, _errs, nil
}

func parseRow(idx imcsv.Indexer, r int, row []string, currentRowFulfillment *RowFulfillment, _errs []error) (_ *RowFulfillment, errs []error) {
	fulfillmentEdCode := idx.GetCell(row, idxFulfillmentEdCode)
	customerName := idx.GetCell(row, idxCustomerName)
	customerPhone := idx.GetCell(row, idxCustomerPhone)
	shippingAddress := idx.GetCell(row, idxShippingAddress)
	province := idx.GetCell(row, idxProvince)
	district := idx.GetCell(row, idxDistrict)
	ward := idx.GetCell(row, idxWard)
	productDescription := idx.GetCell(row, idxProductDescription)
	shippingNote := idx.GetCell(row, idxShippingNote)

	col := idxTotalWeight
	totalWeight, err := imcsv.ParseUint(idx.GetCell(row, idxTotalWeight))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxBasketValue
	basketValue, err := parseMoney(idx.GetCell(row, idxBasketValue))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxIncludeInsurance
	includeInsurance, err := parseBool(idx.GetCell(row, idxIncludeInsurance))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	col = idxCODAmount
	CODAmount, err := parseMoney(idx.GetCell(row, idxCODAmount))
	if err != nil {
		err = imcsv.CellError(idx, r, col, err.Error())
		errs = append(errs, err)
	}

	var rowFulfillment RowFulfillment
	rowFulfillment.RowIndex = r
	rowFulfillment.EdCode = fulfillmentEdCode
	rowFulfillment.CustomerName = customerName
	rowFulfillment.CustomerPhone = customerPhone
	rowFulfillment.ShippingAddress = shippingAddress
	rowFulfillment.Province = province
	rowFulfillment.District = district
	rowFulfillment.Ward = ward
	rowFulfillment.ProductDescription = productDescription
	rowFulfillment.TotalWeight = totalWeight
	rowFulfillment.BasketValue = basketValue
	rowFulfillment.IncludeInsurance = includeInsurance
	rowFulfillment.CODAmount = CODAmount
	rowFulfillment.ShippingNote = shippingNote

	return &rowFulfillment, nil
}

func parseRowsToModels(rows []*RowFulfillment) []*types.FulfillmentResponse {
	var result []*types.FulfillmentResponse
	for _, row := range rows {
		result = append(result, parseRowToModel(row))
	}

	return result
}

func parseRowToModel(row *RowFulfillment) *types.FulfillmentResponse {
	result := &types.FulfillmentResponse{
		EdCode:             row.EdCode,
		CustomerName:       row.CustomerName,
		CustomerPhone:      row.CustomerPhone,
		ShippingAddress:    row.ShippingAddress,
		District:           "",
		DistrictCode:       "",
		Province:           "",
		ProvinceCode:       "",
		Ward:               "",
		WardCode:           "",
		ProductDescription: row.ProductDescription,
		TotalWeight:        row.TotalWeight,
		BasketValue:        row.BasketValue,
		IncludeInsurance:   row.IncludeInsurance,
		CODAmount:          row.CODAmount,
		ShippingNote:       row.ShippingNote,
	}
	shippingAddress := strings.TrimSpace(row.ShippingAddress)
	if len(shippingAddress) != 0 {
		rawLocation := strings.Join([]string{shippingAddress, row.Ward, row.District, row.Province}, " ")
		_location := location.ParseLocation(rawLocation)
		if _location.Province != nil {
			result.Province = _location.Province.Name
			result.ProvinceCode = _location.Province.Code
		}
		if _location.District != nil {
			result.District = _location.District.Name
			result.DistrictCode = _location.District.Code
		}
		if _location.Ward != nil {
			result.Ward = _location.Ward.Name
			result.WardCode = _location.Ward.Code
		}
	}

	return result
}
