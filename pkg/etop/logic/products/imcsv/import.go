package imcsv

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	pbcm "etop.vn/backend/pb/common"
	pbshop "etop.vn/backend/pb/etop/shop"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/imcsv"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
)

func HandleShopImportSampleProducts(c *httpx.Context) error {
	claim := c.Claim.(*claims.ShopClaim)
	userID := c.Session.GetUserID()
	shop := claim.Shop

	// share the same key with HandleShopImportProducts
	key := strconv.FormatInt(shop.ID, 10)

	resp, err := idempgroup.DoAndWrapWithSubkey(key, claim.Token, 30*time.Second, func() (interface{}, error) {
		return handleShopImportSampleProducts(c.Req.Context(), c, shop, userID)
	}, "tạo sản phẩm mẫu")
	if err != nil {
		return err
	}

	respMsg := resp.(*pbshop.ImportProductsResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, claim.Token)
	}
	c.SetResultPb(respMsg)
	return nil
}

func handleShopImportSampleProducts(ctx context.Context, c *httpx.Context, shop *model.Shop, userID int64) (_resp *pbshop.ImportProductsResponse, _err error) {
	if shop.ProductSourceID != 0 {
		// check if shop already imports sample data
		s := productStore(ctx).
			ProductSourceID(shop.ProductSourceID).
			Code("TEST-SP-01")
		products, err := s.ListProducts()
		if err != nil {
			return nil, cm.Error(cm.Internal, "Không thể tạo sản phẩm mẫu", err)
		}

		if len(products) != 0 {
			_resp = &pbshop.ImportProductsResponse{
				ImportErrors: []*pbcm.Error{{Code: "ok", Msg: "Sản phẩm mẫu đã được import"}},
			}
			return
		}
	}

	reader := ioutil.NopCloser(bytes.NewReader(dlShopProductXlsx))
	return handleShopImportProductsFromFile(ctx, c, shop, userID, 0, reader, assetShopProductFilename)
}

func HandleShopImportProducts(c *httpx.Context) error {
	claim := c.Claim.(*claims.ShopClaim)
	userID := c.Session.GetUserID()
	shop := claim.Shop

	// share the same key with HandleShopImportSampleProducts
	key := strconv.FormatInt(shop.ID, 10)

	resp, err := idempgroup.DoAndWrapWithSubkey(key, claim.Token, 30*time.Second, func() (interface{}, error) {
		return handleShopImportProducts(c.Req.Context(), c, shop, userID)
	}, "import đơn hàng")
	if err != nil {
		return err
	}

	respMsg := resp.(*pbshop.ImportProductsResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, claim.Token)
	}
	c.SetResultPb(respMsg)
	return nil
}

func handleShopImportProducts(ctx context.Context, c *httpx.Context, shop *model.Shop, userID int64) (_resp *pbshop.ImportProductsResponse, _err error) {
	mode, fileHeader, err := parseRequest(c)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "can not open file")
	}
	return handleShopImportProductsFromFile(ctx, c, shop, userID, mode, file, fileHeader.Filename)
}

func handleShopImportProductsFromFile(ctx context.Context, c *httpx.Context, shop *model.Shop, userID int64, mode Mode, file io.ReadCloser, filename string) (_resp *pbshop.ImportProductsResponse, _err error) {
	defer file.Close()
	var debugOpts Debug
	if cm.NotProd() {
		var err error
		debugOpts, err = parseDebugHeader(c.Req.Header)
		if err != nil {
			return nil, err
		}
	}

	imp := &Importer{Mode: mode}
	startAt := time.Now()
	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "can not open file")
	}

	// We only store file if the file is valid.
	importID := cm.NewIDWithTag(model.TagImport)
	uploadCmd, err := uploadFile(importID, rawData)
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
			OriginalFile: filename,
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
			attempt.Status = model.S4Negative
			attempt.ErrorType = "panic"
			savedErr := cm.Errorf(cm.Internal, nil, "%v", rerr).
				WithMeta("stack", cm.UnsafeBytesToString(stack))
			attempt.Errors = []*model.Error{model.ToError(savedErr)}

			// respond internal error to client
			_err = cm.Error(cm.Internal, "", nil)

		case _err != nil:
			attempt.Status = model.S4Negative
			attempt.ErrorType = "error"
			err = cm.ToError(_err).WithMetaID("import_id", importID)
			attempt.Errors = []*model.Error{model.ToError(_err)}

		case len(_resp.CellErrors) > 0:
			attempt.Status = model.S4Negative
			attempt.ErrorType = "cell_errors"
			attempt.Errors = pbcm.ErrorsToModel(_resp.CellErrors)
			attempt.NError = len(_resp.CellErrors)

		case len(_resp.ImportErrors) > 0:
			count := pbcm.CountErrors(_resp.ImportErrors)
			if count == 0 {
				attempt.Status = model.S4Positive
				attempt.NCreated = len(_resp.ImportErrors)

			} else {
				attempt.Status = model.S4SuperPos // partially error
				attempt.ErrorType = "import_errors"
				attempt.Errors = pbcm.ErrorsToModel(_resp.ImportErrors)
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
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "invalid file format")
	}

	sheetName, err := validateSheets(excelFile)
	if err != nil {
		return nil, err
	}

	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "no rows")
	}
	imp.Rows = rows

	schema, idx, _errs, err := validateSchema(&rows[0])
	imp.Schema = schema
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(_errs)
	}

	imp.LastRow, _errs, err = validateRows(schema, rows, idx)
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(_errs)
	}

	codeMode, rowProducts, _errs, err := parseRows(schema, idx, imp.Mode, rows, imp.LastRow)
	if err != nil {
		return nil, err
	}
	if len(_errs) > 0 {
		return imp.generateErrorResponse(_errs)
	}

	// create new product source if not exist
	if shop.ProductSourceID == 0 {
		createProductSourceCmd := &catalogmodelx.CreateProductSourceCommand{
			ShopID: shop.ID,
			Name:   shop.Name,
			Type:   catalogmodel.ProductSourceCustom,
		}
		if err := bus.Dispatch(ctx, createProductSourceCmd); err != nil {
			_err = cm.Error(cm.Internal, "", err).
				WithMeta("step", "create product source")
			return
		}
		shop.ProductSourceID = createProductSourceCmd.Result.ID
	}

	// this function expects product source not empty
	requests, _errs := parseRowsToModel(schema, idx, imp.Mode, rowProducts, shop)
	if len(_errs) > 0 {
		return imp.generateErrorResponse(_errs)
	}

	msgs, _errs, _cellErrs, err := loadAndCreateProducts(ctx, schema, idx, imp.Mode, codeMode, shop, rowProducts, requests, debugOpts)
	if err != nil {
		return nil, err
	}
	if len(_cellErrs) > 0 {
		return imp.generateErrorResponse(_cellErrs)
	}

	resp := &pbshop.ImportProductsResponse{
		Data: imp.toSpreadsheetData(idx.indexer),
	}
	importErrors := make([]*pbcm.Error, 0, len(msgs)+len(_errs))
	for _, msg := range msgs {
		importErrors = append(importErrors, &pbcm.Error{
			Code: "ok",
			Msg:  msg,
		})
	}
	for _, err := range _errs {
		importErrors = append(importErrors, pbcm.PbError(err))
	}
	resp.ImportErrors = importErrors
	return resp, nil
}

func validateSheets(file *excelize.File) (sheetName string, err error) {
	sheetName = file.GetSheetName(1)
	if sheetName == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "invalid sheet")
	}

	norm := validate.NormalizeSearchSimple(sheetName)
	if !strings.Contains(norm, "san pham") && !strings.Contains(norm, "sanpham") {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Sheet đầu tiên trong file phải là danh sách sản phẩm cần import. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.").WithMeta("reason", "invalid sheet name")
	}

	return sheetName, nil
}

func validateRows(schema imcsv.Schema, rows [][]string, idx indexes) (lastNonEmptyRow int, errs []error, _ error) {
	if len(rows) > MaxRows {
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "File import quá lớn. Vui lòng kiểm tra lại hoặc liên hệ hotro@etop.vn.")
	}

	for r := 1; r < len(rows); r++ {
		if len(rows[r]) < len(schema) {
			err := imcsv.CellError(idx.indexer, r, -1, "Số cột không đúng cấu trúc yêu cầu.")
			errs = append(errs, err)
			if len(errs) >= MaxCellErrors {
				return
			}
			continue
		}

		imcsv.CleanRow(&rows[r], len(schema))
		if !imcsv.IsRowEmpty(rows[r]) {
			lastNonEmptyRow = r
		}
	}
	return
}

func parseRows(schema imcsv.Schema, idx indexes, mode Mode, rows [][]string, lastRow int) (codeMode CodeMode, products []*RowProduct, _errs []error, _ error) {

	// make sure that either the import file contains product code and variant
	// code for all rows, or leave both column empty
	var init, useCode bool
	for r := 1; r < lastRow; r++ {
		row := rows[r]
		if imcsv.IsRowEmpty(row) {
			continue
		}
		if !init {
			useCode = row[idx.productCode] != ""
			init = true
		}
		productCode, variantCode := row[idx.productCode], row[idx.variantCode]
		switch {
		case useCode && productCode != "" && variantCode != "": // no-op
		case !useCode && productCode == "" && variantCode == "": // no-op
		default:
			err := imcsv.CellError(idx.indexer, r, -1, "Mã sản phẩm và mã phiên bản sản phẩm là không bắt buộc. Nhưng nếu bạn đã cung cấp mã, hãy điền đủ cho tất cả các dòng.")
			_errs = append(_errs, err)
			if len(_errs) >= MaxCellErrors {
				return
			}
		}
	}
	if useCode {
		codeMode = CodeModeUseCode
	} else {
		codeMode = CodeModeUseName
	}
	if len(_errs) != 0 {
		return
	}

	products = make([]*RowProduct, 0, lastRow)
	var rowProduct *RowProduct
	for r := 1; r < len(rows); r++ {
		row := rows[r]
		if imcsv.IsRowEmpty(row) {
			continue
		}
		rowProduct, _errs = parseRow(schema, idx, mode, r, row, _errs)
		products = append(products, rowProduct)
	}
	return
}

func parseRow(schema imcsv.Schema, idx indexes, mode Mode, r int, row []string, _errs []error) (*RowProduct, []error) {
	var col int
	var err error
	rowProduct := &RowProduct{
		RowIndex:      r,
		Category:      [3]string{},
		Collections:   nil,
		ProductCode:   row[idx.productCode],
		VariantCode:   row[idx.variantCode],
		ProductName:   row[idx.productName],
		Attributes:    nil,
		ListPrice:     0,
		CostPrice:     0,
		QuantityAvail: 0,
		Unit:          row[idx.unit],
		ImageURLs:     nil,
		Weight:        0,
		Description:   row[idx.description],
	}

	col = idx.category
	if v := row[col]; v != "" {
		rowProduct.Category, err = parseCategory(v)
		if err != nil {
			err = imcsv.CellError(idx.indexer, r, col, err.Error())
			_errs = append(_errs, err)
		}
	}

	col = idx.collections
	if v := row[col]; v != "" {
		rowProduct.Collections, err = parseCollections(v)
		if err != nil {
			err = imcsv.CellError(idx.indexer, r, col, err.Error())
			_errs = append(_errs, err)
		}
	}

	col = idx.attributes
	if v := row[col]; v != "" {
		rowProduct.Attributes, err = parseAttributes(v)
		if err != nil {
			err = imcsv.CellError(idx.indexer, r, col, err.Error())
			_errs = append(_errs, err)
		}
	}

	col = idx.listPrice
	rowProduct.ListPrice, err = imcsv.ParseUint(row[col])
	if err != nil {
		err = imcsv.CellError(idx.indexer, r, col, err.Error())
		_errs = append(_errs, err)
	}

	col = idx.costPrice
	rowProduct.CostPrice, err = imcsv.ParseUint(row[col])
	if err != nil {
		err = imcsv.CellError(idx.indexer, r, col, err.Error())
		_errs = append(_errs, err)
	}

	col = idx.quantityAvail
	rowProduct.QuantityAvail, err = imcsv.ParseUint(row[col])
	if err != nil {
		err = imcsv.CellError(idx.indexer, r, col, err.Error())
		_errs = append(_errs, err)
	}

	col = idx.images
	if v := row[col]; v != "" {
		rowProduct.ImageURLs, err = parseImageUrls(v)
		if err != nil {
			err = imcsv.CellError(idx.indexer, r, col, err.Error())
			_errs = append(_errs, err)
		}
	}

	col = idx.weight
	if v := row[col]; v != "" {
		rowProduct.Weight, err = parseWeightAsGram(schema, v)
		if err != nil {
			err = imcsv.CellError(idx.indexer, r, col, err.Error())
			_errs = append(_errs, err)
		}
	}

	return rowProduct, _errs
}

func parseCategory(v string) (res [3]string, err error) {
	parts := strings.Split(v, ">>")
	if len(parts) > 3 {
		err = fmt.Errorf("Danh mục chỉ hỗ trợ tối đa 3 cấp")
		return
	}
	// revert it
	for i := 0; i < len(parts)/2; i++ {
		j := len(parts) - 1 - i
		parts[i], parts[j] = parts[j], parts[i]
	}
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		if parts[i] == "" {
			err = fmt.Errorf("Danh mục không hợp lệ.")
			return
		}
		res[i] = parts[i]
	}
	return res, nil
}

func parseCollections(v string) ([]string, error) {
	parts := split(v)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts, nil
}

func parseAttributes(v string) ([]catalogmodel.ProductAttribute, error) {
	items := split(v)
	if len(items) >= 7 {
		return nil, fmt.Errorf(`Quá nhiều thuộc tính.`)
	}
	res := make([]catalogmodel.ProductAttribute, len(items))
	for i, p := range items {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf(`Thuộc tính không hợp lệ (%v). Cần sử dụng cấu trúc "Tên thuộc tính: giá trị thuộc tính".`, p)
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return nil, fmt.Errorf(`Thuộc tính không hợp lệ (%v). Cần sử dụng cấu trúc "Tên thuộc tính: giá trị thuộc tính".`, p)
		}
		res[i] = catalogmodel.ProductAttribute{Name: key, Value: value}
	}
	return res, nil
}

func parseImageUrls(v string) ([]string, error) {
	parts := split(v)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		if !validate.URL(parts[i]) {
			return nil, fmt.Errorf("Link hình ảnh không hợp lệ (%v).", parts[i])
		}
	}
	return parts, nil
}

func parseWeightAsGram(schema imcsv.Schema, v string) (int, error) {
	unit := ""
	vv := v
	switch {
	case strings.HasSuffix(v, "kg") || strings.HasSuffix(v, "Kg"):
		unit = "kg"
		vv = strings.TrimSpace(v[:len(v)-2])
	case strings.HasSuffix(v, "g") || strings.HasSuffix(v, "G"):
		unit = "g"
		vv = strings.TrimSpace(v[:len(v)-1])
	case strings.HasSuffix(v, "gr") || strings.HasSuffix(v, "Gr"):
		unit = "g"
		vv = strings.TrimSpace(v[:len(v)-2])
	case schema[0].Name == schemaV0[0].Name: // TODO: refactor to use schema.Translator
		unit = "g"
	default:
		unit = "kg"
	}

	if unit == "g" {
		return imcsv.ParseUint(vv)
	}

	f, err := imcsv.ParseFloat(v)
	if err != nil {
		return 0, err
	}
	return int(f * 1000), nil
}

func split(v string) []string {
	if strings.Contains(v, "|") {
		return strings.Split(v, "|")
	}
	return strings.Split(v, ",")
}

func parseRowsToModel(schema imcsv.Schema, idx indexes, mode Mode, rowProducts []*RowProduct, shop *model.Shop) (requests []*pbshop.CreateVariantRequest, _errs []error,
) {
	now := time.Now()
	requests = make([]*pbshop.CreateVariantRequest, len(rowProducts))

	for i, rowProduct := range rowProducts {
		errs := rowProduct.Validate(schema, idx, mode)
		if len(errs) > 0 {
			_errs = append(_errs, errs...)
			if len(_errs) > MaxCellErrors {
				return
			}
			continue
		}
		requests[i] = parseRowToModel(rowProduct, shop.ProductSourceID, now)
	}
	return
}

func parseRowToModel(rowProduct *RowProduct, productSourceID int64, now time.Time) *pbshop.CreateVariantRequest {
	return &pbshop.CreateVariantRequest{
		ProductSourceId:   productSourceID, // it may be empty, will be filled later
		ProductId:         0,               // will be filled later
		ProductName:       rowProduct.ProductName,
		Name:              variantNameFromAttributes(rowProduct.Attributes),
		Description:       rowProduct.Description,
		ShortDesc:         "",
		DescHtml:          "",
		ImageUrls:         rowProduct.ImageURLs,
		Tags:              nil,
		Status:            0,
		ListPrice:         int32(rowProduct.ListPrice),
		CostPrice:         int32(rowProduct.CostPrice),
		Sku:               rowProduct.VariantCode,
		Code:              rowProduct.ProductCode,
		QuantityAvailable: int32(rowProduct.QuantityAvail),
		QuantityOnHand:    0,
		QuantityReserved:  0,
		Attributes:        attributesToModel(rowProduct.Attributes),
		Unit:              rowProduct.Unit,
	}
}

func variantNameFromAttributes(attrs []catalogmodel.ProductAttribute) string {
	if len(attrs) == 0 {
		return ""
	}
	var s strings.Builder
	for i, attr := range attrs {
		if i > 0 {
			s.WriteString(" ")
		}
		s.WriteString(attr.Value)
	}
	return s.String()
}

func attributesToModel(attrs []catalogmodel.ProductAttribute) []*pbshop.Attribute {
	if len(attrs) == 0 {
		return nil
	}
	res := make([]*pbshop.Attribute, len(attrs))
	for i, attr := range attrs {
		res[i] = &pbshop.Attribute{
			Name:  attr.Name,
			Value: attr.Value,
		}
	}
	return res
}
