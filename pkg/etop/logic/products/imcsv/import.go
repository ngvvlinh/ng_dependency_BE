package imcsv

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"runtime/debug"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	apishop "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/api/top/types/etc/account_tag"
	"o.o/api/top/types/etc/status4"
	com "o.o/backend/com/main"
	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/etop/upload"
)

var idempgroup *idemp.RedisGroup

const PrefixIdemp = "IdempImportProduct"

type Import struct {
	uploader         *upload.Uploader
	shopProductStore catalogsqlstore.ShopProductStoreFactory
	shopVariantStore catalogsqlstore.ShopVariantStoreFactory

	ExportAttemptStore sqlstore.ExportAttemptStoreInterface
	CategoryStore      sqlstore.CategoryStoreInterface
	ShopStore          sqlstore.ShopStoreInterface
}

func New(
	rd redis.Store,
	ul *upload.Uploader,
	db com.MainDB,
	ExportAttemptStore sqlstore.ExportAttemptStoreInterface,
	CategoryStore sqlstore.CategoryStoreInterface,
	ShopStore sqlstore.ShopStoreInterface,
) (*Import, func()) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60) // 5 minutes
	im := &Import{
		uploader:           ul,
		shopProductStore:   catalogsqlstore.NewShopProductStore(db),
		shopVariantStore:   catalogsqlstore.NewShopVariantStore(db),
		ExportAttemptStore: ExportAttemptStore,
		CategoryStore:      CategoryStore,
		ShopStore:          ShopStore,
	}
	if ul != nil {
		im.uploader = ul
		ul.ExpectDir(model.ImportTypeShopProduct.String())
	}
	return im, idempgroup.Shutdown
}

func (im *Import) HandleShopImportSampleProducts(c *httpx.Context) error {
	shop, user, token := c.SS.Shop(), c.SS.User(), c.SS.Claim().Token
	// share the same key with HandleShopImportProducts
	key := shop.ID.String()
	resp, _, err := idempgroup.DoAndWrapWithSubkey(c.Context(), key, token, 30*time.Second, func() (interface{}, error) {
		return im.handleShopImportSampleProducts(c.Req.Context(), c, shop, user)
	}, "t???o s???n ph???m m???u")
	if err != nil {
		return err
	}

	respMsg := resp.(*apishop.ImportProductsResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, token)
	}
	c.SetResult(respMsg)
	return nil
}

func (im *Import) handleShopImportSampleProducts(ctx context.Context, c *httpx.Context, shop *identitymodel.Shop, user *identitymodelx.SignedInUser) (_resp *apishop.ImportProductsResponse, _err error) {
	// check if shop already imports sample data
	s := im.shopProductStore(ctx).
		ShopID(shop.ID).
		Code("TEST-SP-01")
	products, err := s.ListShopProducts()
	if err != nil {
		return nil, cm.Error(cm.Internal, "Kh??ng th??? t???o s???n ph???m m???u", err)
	}

	if len(products) != 0 {
		_resp = &apishop.ImportProductsResponse{
			ImportErrors: []*pbcm.Error{{Code: "ok", Msg: "S???n ph???m m???u ???? ???????c import"}},
		}
		return
	}

	reader := ioutil.NopCloser(bytes.NewReader(dlShopProductXlsx))
	return im.handleShopImportProductsFromFile(ctx, c, shop, user, 0, reader, assetShopProductFilename)
}

func (im *Import) HandleShopImportProducts(c *httpx.Context) error {
	claim, shop, user := c.SS.Claim(), c.SS.Shop(), c.SS.User()
	// share the same key with HandleShopImportSampleProducts
	key := shop.ID.String()

	resp, _, err := idempgroup.DoAndWrapWithSubkey(c.Context(), key, claim.Token, 30*time.Second, func() (interface{}, error) {
		return im.handleShopImportProducts(c.Req.Context(), c, shop, user)
	}, "import ????n h??ng")
	if err != nil {
		return err
	}

	respMsg := resp.(*apishop.ImportProductsResponse)
	if len(respMsg.CellErrors) > 0 {
		// Allow re-uploading immediately after error
		idempgroup.ReleaseKey(key, claim.Token)
	}
	c.SetResult(respMsg)
	return nil
}

func (im *Import) handleShopImportProducts(ctx context.Context, c *httpx.Context, shop *identitymodel.Shop, user *identitymodelx.SignedInUser) (_resp *apishop.ImportProductsResponse, _err error) {
	mode, fileHeader, err := parseRequest(c)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Kh??ng th??? ?????c ???????c file. Vui l??ng ki???m tra l???i ho???c li??n h??? %v.", wl.X(ctx).CSEmail).WithMeta("reason", "can not open file")
	}
	return im.handleShopImportProductsFromFile(ctx, c, shop, user, mode, file, fileHeader.Filename)
}

func (im *Import) handleShopImportProductsFromFile(ctx context.Context, c *httpx.Context, shop *identitymodel.Shop, user *identitymodelx.SignedInUser, mode Mode, file io.ReadCloser, filename string) (_resp *apishop.ImportProductsResponse, _err error) {
	defer file.Close()
	var debugOpts Debug
	if cmenv.NotProd() {
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
		return nil, cm.Errorf(cm.InvalidArgument, err, "Kh??ng th??? ?????c ???????c file. Vui l??ng ki???m tra l???i ho???c li??n h??? %v.", wl.X(ctx).CSEmail).WithMeta("reason", "can not open file")
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
			UserID:       user.ID,
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
			err = cm.ToError(_err).WithMetaID("import_id", importID)
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
		if err = im.ExportAttemptStore.CreateImportAttempt(ctx, createAttemptCmd); err != nil {
			_err = err
		}
	}()

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Kh??ng th??? ?????c ???????c file. Vui l??ng ki???m tra l???i ho???c li??n h??? %v.", wl.X(ctx).CSEmail).WithMeta("reason", "invalid file format")
	}

	sheetName, err := validateSheets(excelFile)
	if err != nil {
		return nil, err
	}

	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File kh??ng c?? n???i dung. Vui l??ng t???i l???i file import ho???c li??n h??? %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}
	imp.Rows = rows

	schema, idx, _errs, err := validateSchema(c.Context(), &rows[0])
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

	stocktakeId, msgs, _errs, _cellErrs, err := im.loadAndCreateProducts(ctx, c.Session, schema, idx, imp.Mode, codeMode, shop, rowProducts, debugOpts, user)
	if err != nil {
		return nil, err
	}
	if len(_cellErrs) > 0 {
		return imp.generateErrorResponse(_cellErrs)
	}

	resp := &apishop.ImportProductsResponse{
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
		importErrors = append(importErrors, cmapi.PbError(err))
	}
	resp.ImportErrors = importErrors
	resp.StocktakeID = stocktakeId
	return resp, nil
}

func validateSheets(file *excelize.File) (sheetName string, err error) {
	sheetName = file.GetSheetName(1)
	if sheetName == "" {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Kh??ng th??? ?????c ???????c file.").WithMeta("reason", "invalid sheet")
	}

	norm := validate.NormalizeSearchSimple(sheetName)
	if !strings.Contains(norm, "san pham") && !strings.Contains(norm, "sanpham") {
		return "", cm.Errorf(cm.InvalidArgument, nil, "Sheet ?????u ti??n trong file ph???i l?? danh s??ch s???n ph???m c???n import.").WithMeta("reason", "invalid sheet name")
	}

	return sheetName, nil
}

func validateRows(schema imcsv.Schema, rows [][]string, idx indexes) (lastNonEmptyRow int, errs []error, _ error) {
	if len(rows) > MaxRows {
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "File import qu?? l???n.")
	}

	for r := 1; r < len(rows); r++ {
		if len(rows[r]) < len(schema) {
			err := imcsv.CellError(idx.indexer, r, -1, "S??? c???t kh??ng ????ng c???u tr??c y??u c???u.")
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
			err := imcsv.CellError(idx.indexer, r, -1, "M?? s???n ph???m v?? m?? phi??n b???n s???n ph???m l?? kh??ng b???t bu???c. Nh??ng n???u b???n ???? cung c???p m??, h??y ??i???n ????? cho t???t c??? c??c d??ng.")
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

	if err == nil {
		errs := rowProduct.Validate(schema, idx, mode)
		_errs = append(_errs, errs...)
	}
	return rowProduct, _errs
}

func parseCategory(v string) (res [3]string, err error) {
	parts := strings.Split(v, ">>")
	if len(parts) > 3 {
		err = fmt.Errorf("Danh m???c ch??? h??? tr??? t???i ??a 3 c???p")
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
			err = fmt.Errorf("Danh m???c kh??ng h???p l???.")
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

func parseAttributes(v string) ([]*catalogmodel.ProductAttribute, error) {
	items := split(v)
	if len(items) >= 7 {
		return nil, fmt.Errorf(`Qu?? nhi???u thu???c t??nh.`)
	}
	res := make([]*catalogmodel.ProductAttribute, len(items))
	for i, p := range items {
		parts := strings.Split(p, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf(`Thu???c t??nh kh??ng h???p l??? (%v). C???n s??? d???ng c???u tr??c "T??n thu???c t??nh: gi?? tr??? thu???c t??nh".`, p)
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return nil, fmt.Errorf(`Thu???c t??nh kh??ng h???p l??? (%v). C???n s??? d???ng c???u tr??c "T??n thu???c t??nh: gi?? tr??? thu???c t??nh".`, p)
		}
		res[i] = &catalogmodel.ProductAttribute{Name: key, Value: value}
	}
	return res, nil
}

func parseImageUrls(v string) ([]string, error) {
	parts := split(v)
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
		if !validate.URL(parts[i]) {
			return nil, fmt.Errorf("Link h??nh ???nh kh??ng h???p l??? (%v).", parts[i])
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

func rowToCreateVariant(row *RowProduct, now time.Time) *apishop.CreateVariantRequest {
	return &apishop.CreateVariantRequest{
		Code:        row.VariantCode,
		Name:        variantNameFromAttributes(row.Attributes),
		Attributes:  convert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(row.Attributes),
		ProductId:   0, // will be filled later
		Note:        "",
		Description: row.Description,
		ShortDesc:   "",
		DescHtml:    "",
		ImageUrls:   row.ImageURLs,
		CostPrice:   row.CostPrice,
		ListPrice:   row.ListPrice,
		RetailPrice: row.ListPrice,
	}
}

func rowToCreateProduct(row *RowProduct, now time.Time) *apishop.CreateProductRequest {
	return &apishop.CreateProductRequest{
		Code:        row.ProductCode,
		Name:        row.ProductName,
		Unit:        row.Unit,
		Note:        "",
		Description: "",
		ShortDesc:   "",
		DescHtml:    "",
		ImageUrls:   row.ImageURLs,
		CostPrice:   row.CostPrice,
		ListPrice:   row.ListPrice,
		RetailPrice: 0,
	}
}

func variantNameFromAttributes(attrs []*catalogmodel.ProductAttribute) string {
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
