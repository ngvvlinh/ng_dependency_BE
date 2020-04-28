package imcsv

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"mime/multipart"
	"net/http"

	"o.o/api/top/int/shop"
	pbsheet "o.o/api/top/int/types/spreadsheet"
	catalogsqlstore "o.o/backend/com/main/catalog/sqlstore"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/idemp"
	cmservice "o.o/backend/pkg/common/apifw/service"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup
var uploader *upload.Uploader
var shopProductStore catalogsqlstore.ShopProductStoreFactory
var shopVariantStore catalogsqlstore.ShopVariantStoreFactory

const PrefixIdemp = "IdempImportProduct"

func Init(sd cmservice.Shutdowner, rd redis.Store, ul *upload.Uploader, db *cmsql.Database) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60) // 5 minutes
	sd.Register(idempgroup.Shutdown)
	shopProductStore = catalogsqlstore.NewShopProductStore(db)
	shopVariantStore = catalogsqlstore.NewShopVariantStore(db)

	if ul != nil {
		uploader = ul
		ul.ExpectDir(model.ImportTypeShopProduct.String())
	}
}

const (
	MaxFilesize   = 2 * 1024 * 1024 // 2 MB
	MaxRows       = 1200
	MaxCellErrors = 20
	MaxEmptyRows  = 20
)

// Mode has not been implemented yet
type Mode int

func ParseMode(v string) (Mode, error) {
	return 0, nil
}

type CodeMode int

const (
	CodeModeUseCode CodeMode = 1
	CodeModeUseName CodeMode = 2
)

type Importer struct {
	Schema  imcsv.Schema
	Mode    Mode
	Rows    [][]string
	LastRow int
}

func parseRequest(c *httpx.Context) (Mode, *multipart.FileHeader, error) {
	// Limit the max file size
	c.Req.Body = http.MaxBytesReader(c.Resp, c.Req.Body, MaxFilesize)

	mode, err := ParseMode(c.Req.FormValue("mode"))
	if err != nil {
		return 0, nil, err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
		// continue
	default:
		return 0, nil, cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}
	return mode, files[0], nil
}

func (imp *Importer) generateErrorResponse(errs []error) (*shop.ImportProductsResponse, error) {
	resp := &shop.ImportProductsResponse{
		Data:       imp.toSpreadsheetData(imcsv.Indexer{}),
		CellErrors: cmapi.PbErrors(errs),
	}
	return resp, nil
}

func (imp *Importer) toSpreadsheetData(idx imcsv.Indexer) *pbsheet.SpreadsheetData {
	return imcsv.ToSpreadsheetData(imp.Schema, idx, imp.Rows, imp.LastRow)
}

func uploadFile(id dot.ID, data []byte) (*upload.StoreFileCommand, error) {
	fileName := fmt.Sprintf("%v.xlsx", id)
	uploadCmd := &upload.StoreFileCommand{
		UploadType: model.ImportTypeShopProduct.String(),
		FileName:   fileName,
		Data:       data,
	}

	if uploader == nil {
		ll.Warn("Disabled uploading file")
		return uploadCmd, nil
	}
	return uploadCmd, uploader.StoreFile(uploadCmd)
}

type Debug struct {
	FailPercent int `json:"fail_percent"`
}

func parseDebugHeader(h http.Header) (debug Debug, err error) {
	debugHeader := h.Get("debug")
	if debugHeader == "" {
		return
	}

	if err := jsonx.Unmarshal([]byte(debugHeader), &debug); err != nil {
		return Debug{}, cm.Error(cm.InvalidArgument, "Can not read debug header", err)
	}
	if debug.FailPercent < 0 || debug.FailPercent > 100 {
		return Debug{}, cm.Error(cm.InvalidArgument, "Invalid fail_percent", nil)
	}
	return
}

func isRandomFail(percent int) bool {
	var b [2]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	u16 := binary.LittleEndian.Uint16(b[:])
	return u16%100 < uint16(percent)
}
